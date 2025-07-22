#!/usr/bin/env python3
"""
Final MeMeMe TCG Card Scraper
Downloads all card data and images from https://mememe-tcg.com/cardlist
Based on actual HTML structure analysis
"""

import requests
from bs4 import BeautifulSoup
import json
import os
import re
import time
from urllib.parse import urljoin
from typing import Dict, List, Optional
import logging

# Set up logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class MeMeMeCardScraperFinal:
    def __init__(self):
        self.base_url = "https://mememe-tcg.com"
        self.cardlist_url = f"{self.base_url}/cardlist"
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
        })
        
    def parse_cost_info(self, modal_element) -> Dict:
        """Parse cost information from modal element"""
        cost_data = {
            "total": 0,
            "red": 0,
            "blue": 0,
            "yellow": 0,
            "green": 0,
            "colorless": 0
        }
        
        # Get total cost
        total_cost_elem = modal_element.find('span', class_='p-modalCost__totalNum')
        if total_cost_elem:
            cost_data["total"] = int(total_cost_elem.get_text(strip=True))
        
        # Count color symbols from cost icons
        cost_icons = modal_element.find_all('div', class_='p-modalCost__icon')
        for icon in cost_icons:
            img = icon.find('img')
            if img:
                src = img.get('src', '')
                alt = img.get('alt', '')
                
                if 'cost-red' in src or '赤' in alt:
                    cost_data["red"] += 1
                elif 'cost-blue' in src or '青' in alt:
                    cost_data["blue"] += 1
                elif 'cost-yellow' in src or '黄' in alt or '黃' in alt:
                    cost_data["yellow"] += 1
                elif 'cost-green' in src or '緑' in alt:
                    cost_data["green"] += 1
                elif 'cost-mu' in src or '無色' in alt:
                    cost_data["colorless"] += 1
                    
        return cost_data
    
    def parse_card_modal(self, modal_element) -> Optional[Dict]:
        """Parse individual card data from modal element"""
        try:
            card_data = {}
            
            # Card number
            num_elem = modal_element.find('div', class_='p-modalHeadInfo__num')
            if num_elem:
                card_data['number'] = num_elem.get_text(strip=True)
            
            # Card name
            name_elem = modal_element.find('div', class_='p-modalHeadTitle')
            if name_elem:
                card_data['name'] = name_elem.get_text(strip=True)
            
            # Rarity
            rarity_elem = modal_element.find('div', class_='p-modalHeadInfo__rare')
            if rarity_elem:
                card_data['rarity'] = rarity_elem.get_text(strip=True)
            
            # Cost information
            card_data['cost'] = self.parse_cost_info(modal_element)
            
            # Properties
            properties = modal_element.find_all('div', class_='p-modalPropertiesItem')
            for prop in properties:
                tag_elem = prop.find('div', class_='p-modalPropertiesItem__tag')
                value_elem = prop.find('div', class_='p-modalPropertiesItemArea__txt')
                
                if tag_elem and value_elem:
                    tag = tag_elem.get_text(strip=True)
                    value = value_elem.get_text(strip=True)
                    
                    if tag == '色':
                        card_data['color'] = value
                    elif tag == 'カードタイプ':
                        card_data['type'] = value
                    elif tag == '属性':
                        card_data['attribute'] = value
                    elif tag == '感情':
                        card_data['emotion'] = value
                    elif tag == 'パワー':
                        try:
                            card_data['power'] = int(value)
                        except:
                            card_data['power'] = value
                    elif tag == '能力':
                        card_data['ability'] = value
                    elif tag == 'フレーバーテキスト':
                        # Get all text from flavor text area
                        flavor_area = prop.find('div', class_='p-modalPropertiesItemArea')
                        if flavor_area:
                            card_data['flavor_text'] = flavor_area.get_text(separator=' ', strip=True)
            
            # Profile data (height, weight)
            profile = modal_element.find('div', class_='p-modalProfile')
            if profile:
                profile_text = profile.get_text()
                
                # Height
                height_match = re.search(r'身長[：:]?\s*(\d+\.?\d*)\s*cm', profile_text)
                if height_match:
                    card_data['height'] = height_match.group(1)
                
                # Weight
                weight_match = re.search(r'体重[：:]?\s*(\d+\.?\d*)\s*kg', profile_text)
                if weight_match:
                    card_data['weight'] = weight_match.group(1)
            
            # Image URL
            img_elem = modal_element.find('div', class_='p-modalImg')
            if img_elem:
                img = img_elem.find('img')
                if img:
                    src = img.get('src', '')
                    if src:
                        card_data['image_url'] = urljoin(self.base_url, src)
            
            # Determine if promo/parallel
            if 'number' in card_data:
                if '(P)' in card_data['number'] or card_data['number'].endswith('-P'):
                    card_data['is_promo'] = True
                else:
                    card_data['is_promo'] = False
                    
            if 'rarity' in card_data and '-P' in card_data['rarity']:
                card_data['is_parallel'] = True
            else:
                card_data['is_parallel'] = False
            
            return card_data
            
        except Exception as e:
            logger.error(f"Error parsing card modal: {e}")
            return None
    
    def scrape_all_cards(self) -> List[Dict]:
        """Scrape all cards from the cardlist page"""
        logger.info(f"Fetching cardlist from {self.cardlist_url}")
        
        try:
            response = self.session.get(self.cardlist_url)
            response.raise_for_status()
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Find all card modals
            card_modals = soup.find_all('div', class_='p-modal')
            logger.info(f"Found {len(card_modals)} card modals")
            
            all_cards = []
            for i, modal in enumerate(card_modals):
                logger.info(f"Processing card {i+1}/{len(card_modals)}")
                card_data = self.parse_card_modal(modal)
                if card_data:
                    all_cards.append(card_data)
                
            return all_cards
            
        except Exception as e:
            logger.error(f"Error scraping cards: {e}")
            return []
    
    def download_card_image(self, image_url: str, card_number: str, rarity: str) -> bool:
        """Download a card image"""
        try:
            # Create filename based on card number and rarity
            # Clean up the card number for filename
            clean_number = card_number.replace(' ', '').replace('(P)', '-P')
            
            # Determine filename based on whether it's a promo/parallel
            if '-P' in clean_number or rarity.endswith('-P'):
                # For promo/parallel cards, include the suffix
                filename = f"{clean_number}_{rarity}.jpg"
            else:
                # For regular cards
                filename = f"{clean_number}_{rarity}.jpg"
                
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
    
    def save_card_data(self, cards: List[Dict], filename: str = 'all_cards_complete.json'):
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
        logger.info("Starting final card scraping...")
        all_cards = self.scrape_all_cards()
        
        if not all_cards:
            logger.warning("No cards were scraped.")
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
                time.sleep(0.1)  # Be respectful to the server
        
        logger.info("Scraping completed!")
        
        # Generate detailed summary
        summary = {
            "total_cards": len(all_cards),
            "promo_cards": sum(1 for c in all_cards if c.get('is_promo', False)),
            "parallel_cards": sum(1 for c in all_cards if c.get('is_parallel', False)),
            "cards_with_names": sum(1 for c in all_cards if c.get('name')),
            "cards_with_costs": sum(1 for c in all_cards if c.get('cost', {}).get('total', 0) > 0),
            "cards_with_abilities": sum(1 for c in all_cards if c.get('ability')),
            "card_types": list(set(c.get('type', '') for c in all_cards if c.get('type'))),
            "colors": list(set(c.get('color', '') for c in all_cards if c.get('color'))),
            "rarities": list(set(c.get('rarity', '') for c in all_cards if c.get('rarity'))),
            "sample_cards": all_cards[:3] if all_cards else []
        }
        
        with open('data/scraping_summary_final.json', 'w', encoding='utf-8') as f:
            json.dump(summary, f, ensure_ascii=False, indent=2)
        
        logger.info(f"Summary: Total cards: {summary['total_cards']}, "
                   f"Cards with names: {summary['cards_with_names']}, "
                   f"Cards with costs: {summary['cards_with_costs']}")


if __name__ == "__main__":
    scraper = MeMeMeCardScraperFinal()
    scraper.run()