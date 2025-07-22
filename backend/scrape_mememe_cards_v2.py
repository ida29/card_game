#!/usr/bin/env python3
"""
Improved MeMeMe TCG Card Scraper v2
Downloads all card data and images from https://mememe-tcg.com/cardlist
"""

import requests
from bs4 import BeautifulSoup
import json
import os
import re
import time
from urllib.parse import urljoin
from typing import Dict, List, Optional, Tuple
import logging

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class MeMeMeCardScraperV2:
    def __init__(self):
        self.base_url = "https://mememe-tcg.com"
        self.cardlist_url = f"{self.base_url}/cardlist"
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
        })
        
    def save_html_for_debugging(self, html_content: str, filename: str = 'cardlist_debug.html'):
        """Save HTML content for debugging purposes"""
        filepath = os.path.join('data', filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(html_content)
        logger.info(f"Saved HTML for debugging: {filepath}")
        
    def extract_text_between(self, element, start_text: str, end_text: str) -> Optional[str]:
        """Extract text between two markers"""
        text = element.get_text()
        start = text.find(start_text)
        if start == -1:
            return None
        start += len(start_text)
        end = text.find(end_text, start)
        if end == -1:
            return text[start:].strip()
        return text[start:end].strip()
        
    def parse_cost_from_text(self, text: str) -> Dict:
        """Parse cost information from text"""
        cost_data = {
            "total": 0,
            "red": 0,
            "blue": 0,
            "yellow": 0,
            "green": 0,
            "colorless": 0
        }
        
        # Look for cost patterns
        cost_match = re.search(r'コスト[：:]?\s*(\d+)', text)
        if cost_match:
            cost_data["total"] = int(cost_match.group(1))
            
        # Count color symbols
        cost_data["red"] = text.count('赤') + text.count('Red')
        cost_data["blue"] = text.count('青') + text.count('Blue')
        cost_data["yellow"] = text.count('黄') + text.count('Yellow')
        cost_data["green"] = text.count('緑') + text.count('Green')
        
        # Calculate colorless
        if cost_data["total"] > 0:
            color_sum = sum([cost_data["red"], cost_data["blue"], cost_data["yellow"], cost_data["green"]])
            cost_data["colorless"] = max(0, cost_data["total"] - color_sum)
            
        return cost_data
    
    def parse_card_data_v2(self, card_text: str, card_index: int) -> Optional[Dict]:
        """Parse card data from text block"""
        try:
            card_data = {}
            lines = card_text.strip().split('\\n')
            
            # Extract card number (usually first line or contains F-XXX pattern)
            for line in lines:
                number_match = re.search(r'(F-\\d+(?:-P)?(?:\\s*\\(P\\))?)', line)
                if number_match:
                    card_data['number'] = number_match.group(1).strip()
                    break
            
            # Extract card name (usually after card number)
            name_found = False
            for i, line in enumerate(lines):
                if 'number' in card_data and card_data['number'] in line:
                    # Next non-empty line after number is usually the name
                    for j in range(i + 1, len(lines)):
                        if lines[j].strip() and not any(keyword in lines[j] for keyword in ['レアリティ', 'タイプ', 'コスト', '属性']):
                            card_data['name'] = lines[j].strip()
                            name_found = True
                            break
                if name_found:
                    break
            
            # Extract from structured text
            full_text = ' '.join(lines)
            
            # Rarity
            rarity_match = re.search(r'レアリティ[：:]?\\s*([A-Z]+(?:-P)?)', full_text)
            if rarity_match:
                card_data['rarity'] = rarity_match.group(1)
            
            # Color
            color_match = re.search(r'色[：:]?\\s*([赤青黄緑無])', full_text)
            if color_match:
                card_data['color'] = color_match.group(1)
            else:
                # Try alternative patterns
                if '赤' in full_text:
                    card_data['color'] = '赤'
                elif '青' in full_text:
                    card_data['color'] = '青'
                elif '黄' in full_text:
                    card_data['color'] = '黄'
                elif '緑' in full_text:
                    card_data['color'] = '緑'
            
            # Card type
            type_match = re.search(r'タイプ[：:]?\\s*(ふれんど|サポート|フィールド)', full_text)
            if type_match:
                card_data['type'] = type_match.group(1)
            
            # Power
            power_match = re.search(r'パワー[：:]?\\s*(\\d+)', full_text)
            if power_match:
                card_data['power'] = int(power_match.group(1))
            
            # Cost
            card_data['cost'] = self.parse_cost_from_text(full_text)
            
            # Attribute
            attr_match = re.search(r'属性[：:]?\\s*([^\\s\\n]+)', full_text)
            if attr_match:
                card_data['attribute'] = attr_match.group(1)
            
            # Emotion
            emotion_match = re.search(r'感情[：:]?\\s*([^\\s\\n]+)', full_text)
            if emotion_match:
                card_data['emotion'] = emotion_match.group(1)
            
            # Ability/Effect (usually longer text block)
            ability_markers = ['効果', 'アビリティ', '能力']
            for marker in ability_markers:
                if marker in full_text:
                    ability_start = full_text.find(marker)
                    if ability_start != -1:
                        ability_text = full_text[ability_start + len(marker):].strip()
                        # Stop at flavor text or other markers
                        for end_marker in ['フレーバー', '身長', '体重', 'イラスト']:
                            end_pos = ability_text.find(end_marker)
                            if end_pos != -1:
                                ability_text = ability_text[:end_pos].strip()
                        if ability_text and len(ability_text) > 10:
                            card_data['ability'] = ability_text.strip('：:').strip()
                        break
            
            # Flavor text
            flavor_match = re.search(r'フレーバー[：:]?\\s*(.+?)(?=身長|体重|$)', full_text, re.DOTALL)
            if flavor_match:
                card_data['flavor_text'] = flavor_match.group(1).strip()
            
            # Height
            height_match = re.search(r'身長[：:]?\\s*(\\d+\\.?\\d*)\\s*cm', full_text)
            if height_match:
                card_data['height'] = height_match.group(1)
            
            # Weight
            weight_match = re.search(r'体重[：:]?\\s*(\\d+\\.?\\d*)\\s*kg', full_text)
            if weight_match:
                card_data['weight'] = weight_match.group(1)
            
            # Determine if promo/parallel
            if 'number' in card_data:
                if '(P)' in card_data['number'] or card_data['number'].endswith('-P'):
                    card_data['is_promo'] = True
                else:
                    card_data['is_promo'] = False
                    
            if 'rarity' in card_data and card_data['rarity'].endswith('-P'):
                card_data['is_parallel'] = True
            else:
                card_data['is_parallel'] = False
            
            # Image URL (construct based on card number)
            if 'number' in card_data:
                # Clean up card number for URL
                clean_number = card_data['number'].replace(' ', '').replace('(P)', '-P')
                if card_data.get('is_promo'):
                    if 'color' in card_data:
                        color_map = {'赤': 'Red', '青': 'Blue', '黄': 'Yellow', '緑': 'Green'}
                        color_eng = color_map.get(card_data['color'], '')
                        if color_eng:
                            card_data['image_url'] = f"{self.base_url}/assets/images/card/{clean_number}_{card_data['rarity']}_{color_eng}.jpg"
                        else:
                            card_data['image_url'] = f"{self.base_url}/assets/images/card/{clean_number}.jpg"
                    else:
                        card_data['image_url'] = f"{self.base_url}/assets/images/card/{clean_number}.jpg"
                else:
                    card_data['image_url'] = f"{self.base_url}/assets/images/card/{clean_number}.jpg"
            
            return card_data
            
        except Exception as e:
            logger.error(f"Error parsing card at index {card_index}: {e}")
            return None
    
    def scrape_all_cards_v2(self) -> List[Dict]:
        """Scrape all cards using text-based parsing"""
        logger.info(f"Fetching cardlist from {self.cardlist_url}")
        
        try:
            response = self.session.get(self.cardlist_url)
            response.raise_for_status()
            
            # Save HTML for debugging
            self.save_html_for_debugging(response.text)
            
            # Get text content
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Remove script and style elements
            for script in soup(["script", "style"]):
                script.decompose()
                
            text = soup.get_text()
            
            # Split by card number pattern
            card_sections = re.split(r'(?=F-\\d{3})', text)
            card_sections = [s.strip() for s in card_sections if s.strip() and re.search(r'F-\\d{3}', s)]
            
            logger.info(f"Found {len(card_sections)} potential card sections")
            
            all_cards = []
            for i, section in enumerate(card_sections):
                if len(section) > 50:  # Filter out too short sections
                    logger.info(f"Processing card section {i+1}/{len(card_sections)}")
                    card_data = self.parse_card_data_v2(section, i)
                    if card_data and 'number' in card_data:
                        all_cards.append(card_data)
                    time.sleep(0.05)
                
            # Deduplicate by card number
            seen = set()
            unique_cards = []
            for card in all_cards:
                if card['number'] not in seen:
                    seen.add(card['number'])
                    unique_cards.append(card)
                    
            return unique_cards
            
        except Exception as e:
            logger.error(f"Error scraping cards: {e}")
            return []
    
    def download_card_image(self, image_url: str, card_number: str, rarity: str) -> bool:
        """Download a card image"""
        try:
            # Clean filename
            filename = f"{card_number.replace(' ', '').replace('(P)', '-P')}_{rarity}.jpg"
            filepath = os.path.join('data', 'card_images', filename)
            
            # Skip if already downloaded
            if os.path.exists(filepath):
                logger.info(f"Image already exists: {filename}")
                return True
            
            # Download image
            response = self.session.get(image_url)
            response.raise_for_status()
            
            with open(filepath, 'wb') as f:
                f.write(response.content)
            
            logger.info(f"Downloaded: {filename}")
            return True
            
        except Exception as e:
            logger.error(f"Error downloading image {image_url}: {e}")
            return False
    
    def save_card_data(self, cards: List[Dict], filename: str = 'all_cards_v2.json'):
        """Save card data to JSON file"""
        filepath = os.path.join('data', filename)
        with open(filepath, 'w', encoding='utf-8') as f:
            json.dump(cards, f, ensure_ascii=False, indent=2)
        logger.info(f"Saved {len(cards)} cards to {filepath}")
    
    def run(self):
        """Main execution method"""
        # Create necessary directories
        os.makedirs('data', exist_ok=True)
        os.makedirs(os.path.join('data', 'card_images'), exist_ok=True)
        
        # Scrape all cards
        logger.info("Starting card scraping v2...")
        all_cards = self.scrape_all_cards_v2()
        
        if not all_cards:
            logger.warning("No cards were scraped using v2 method.")
            return
        
        # Save card data
        self.save_card_data(all_cards)
        
        # Download images
        logger.info("Starting image downloads...")
        for card in all_cards:
            if 'image_url' in card and 'number' in card and 'rarity' in card:
                self.download_card_image(
                    card['image_url'],
                    card['number'],
                    card['rarity']
                )
                time.sleep(0.2)
        
        logger.info("Scraping completed!")
        
        # Generate summary
        summary = {
            "total_cards": len(all_cards),
            "promo_cards": sum(1 for c in all_cards if c.get('is_promo', False)),
            "parallel_cards": sum(1 for c in all_cards if c.get('is_parallel', False)),
            "cards_with_names": sum(1 for c in all_cards if c.get('name')),
            "cards_with_costs": sum(1 for c in all_cards if c.get('cost', {}).get('total', 0) > 0),
            "cards_with_abilities": sum(1 for c in all_cards if c.get('ability')),
            "card_types": list(set(c.get('type', '') for c in all_cards if c.get('type'))),
            "colors": list(set(c.get('color', '') for c in all_cards if c.get('color'))),
            "rarities": list(set(c.get('rarity', '') for c in all_cards if c.get('rarity')))
        }
        
        with open('data/scraping_summary_v2.json', 'w', encoding='utf-8') as f:
            json.dump(summary, f, ensure_ascii=False, indent=2)
        
        logger.info(f"Summary: {summary}")


if __name__ == "__main__":
    scraper = MeMeMeCardScraperV2()
    scraper.run()