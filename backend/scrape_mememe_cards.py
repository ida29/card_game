#!/usr/bin/env python3
"""
Comprehensive MeMeMe TCG Card Scraper
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

class MeMeMeCardScraper:
    def __init__(self):
        self.base_url = "https://mememe-tcg.com"
        self.cardlist_url = f"{self.base_url}/cardlist"
        self.session = requests.Session()
        self.session.headers.update({
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36'
        })
        
    def parse_cost_info(self, cost_element) -> Dict:
        """Parse cost information from HTML element"""
        cost_data = {
            "total": 0,
            "red": 0,
            "blue": 0,
            "yellow": 0,
            "green": 0,
            "colorless": 0
        }
        
        if not cost_element:
            return cost_data
            
        # Extract total cost from text
        cost_text = cost_element.get_text(strip=True)
        total_match = re.search(r'(\d+)', cost_text)
        if total_match:
            cost_data["total"] = int(total_match.group(1))
        
        # Count color symbols from images
        cost_images = cost_element.find_all('img')
        for img in cost_images:
            src = img.get('src', '')
            alt = img.get('alt', '')
            
            if 'red' in src.lower() or '赤' in alt:
                cost_data["red"] += 1
            elif 'blue' in src.lower() or '青' in alt:
                cost_data["blue"] += 1
            elif 'yellow' in src.lower() or '黄' in alt:
                cost_data["yellow"] += 1
            elif 'green' in src.lower() or '緑' in alt:
                cost_data["green"] += 1
            elif 'colorless' in src.lower() or '無色' in alt or 'gray' in src.lower():
                cost_data["colorless"] += 1
                
        return cost_data
    
    def parse_card_data(self, card_element) -> Optional[Dict]:
        """Parse individual card data from HTML element"""
        try:
            card_data = {}
            
            # Extract card number
            card_num_elem = card_element.find(class_='card-number')
            if card_num_elem:
                card_data['number'] = card_num_elem.get_text(strip=True)
            else:
                # Try alternative methods
                card_id_elem = card_element.find(['span', 'div'], string=re.compile(r'F-\d+'))
                if card_id_elem:
                    card_data['number'] = card_id_elem.get_text(strip=True)
            
            # Extract card name
            name_elem = card_element.find(class_='card-name')
            if name_elem:
                card_data['name'] = name_elem.get_text(strip=True)
            else:
                # Try alternative methods
                name_elem = card_element.find(['h3', 'h4', 'span'], class_=re.compile('name'))
                if name_elem:
                    card_data['name'] = name_elem.get_text(strip=True)
            
            # Extract rarity
            rarity_elem = card_element.find(class_='card-rarity')
            if rarity_elem:
                card_data['rarity'] = rarity_elem.get_text(strip=True)
            else:
                # Look for rarity in text
                rarity_match = re.search(r'\b(C|U|R|SR|SEC)(-P)?\b', str(card_element))
                if rarity_match:
                    card_data['rarity'] = rarity_match.group(0)
            
            # Extract color
            color_elem = card_element.find(class_='card-color')
            if color_elem:
                card_data['color'] = color_elem.get_text(strip=True)
            else:
                # Look for color keywords
                color_text = str(card_element)
                if '赤' in color_text:
                    card_data['color'] = '赤'
                elif '青' in color_text:
                    card_data['color'] = '青'
                elif '黄' in color_text:
                    card_data['color'] = '黄'
                elif '緑' in color_text:
                    card_data['color'] = '緑'
            
            # Extract card type
            type_elem = card_element.find(class_='card-type')
            if type_elem:
                card_data['type'] = type_elem.get_text(strip=True)
            else:
                # Look for type keywords
                type_text = str(card_element)
                if 'ふれんど' in type_text:
                    card_data['type'] = 'ふれんど'
                elif 'サポート' in type_text:
                    card_data['type'] = 'サポート'
                elif 'フィールド' in type_text:
                    card_data['type'] = 'フィールド'
            
            # Extract power
            power_elem = card_element.find(class_='card-power')
            if power_elem:
                power_text = power_elem.get_text(strip=True)
                power_match = re.search(r'(\d+)', power_text)
                if power_match:
                    card_data['power'] = int(power_match.group(1))
            
            # Extract cost information
            cost_elem = card_element.find(class_='card-cost')
            if cost_elem:
                card_data['cost'] = self.parse_cost_info(cost_elem)
            
            # Extract attribute
            attr_elem = card_element.find(class_='card-attribute')
            if attr_elem:
                card_data['attribute'] = attr_elem.get_text(strip=True)
            
            # Extract emotion
            emotion_elem = card_element.find(class_='card-emotion')
            if emotion_elem:
                card_data['emotion'] = emotion_elem.get_text(strip=True)
            
            # Extract ability/effect text
            ability_elem = card_element.find(class_='card-ability')
            if ability_elem:
                card_data['ability'] = ability_elem.get_text(strip=True)
            else:
                ability_elem = card_element.find(class_='card-text')
                if ability_elem:
                    card_data['ability'] = ability_elem.get_text(strip=True)
            
            # Extract flavor text
            flavor_elem = card_element.find(class_='card-flavor')
            if flavor_elem:
                card_data['flavor_text'] = flavor_elem.get_text(strip=True)
            
            # Extract height/weight if available
            height_elem = card_element.find(string=re.compile(r'身長'))
            if height_elem:
                height_match = re.search(r'(\d+\.?\d*)cm', height_elem)
                if height_match:
                    card_data['height'] = height_match.group(1)
            
            weight_elem = card_element.find(string=re.compile(r'体重'))
            if weight_elem:
                weight_match = re.search(r'(\d+\.?\d*)kg', weight_elem)
                if weight_match:
                    card_data['weight'] = weight_match.group(1)
            
            # Extract image URL
            img_elem = card_element.find('img')
            if img_elem:
                img_src = img_elem.get('src', '')
                if img_src:
                    card_data['image_url'] = urljoin(self.base_url, img_src)
            
            # Determine if it's a promo or parallel card
            if 'number' in card_data:
                if '(P)' in card_data['number'] or card_data['number'].endswith('-P'):
                    card_data['is_promo'] = True
                else:
                    card_data['is_promo'] = False
                    
            if 'rarity' in card_data and card_data['rarity'].endswith('-P'):
                card_data['is_parallel'] = True
            else:
                card_data['is_parallel'] = False
            
            return card_data
            
        except Exception as e:
            logger.error(f"Error parsing card: {e}")
            return None
    
    def scrape_all_cards(self) -> List[Dict]:
        """Scrape all cards from the cardlist page"""
        logger.info(f"Fetching cardlist from {self.cardlist_url}")
        
        try:
            response = self.session.get(self.cardlist_url)
            response.raise_for_status()
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Find all card containers - try multiple possible selectors
            card_containers = soup.find_all(class_='card-item')
            if not card_containers:
                card_containers = soup.find_all(class_='card')
            if not card_containers:
                card_containers = soup.find_all('div', id=re.compile(r'F-\d+'))
            if not card_containers:
                # Try to find cards by looking for card number patterns
                card_containers = []
                for elem in soup.find_all(['div', 'section', 'article']):
                    if elem.find(string=re.compile(r'F-\d+')):
                        card_containers.append(elem)
            
            logger.info(f"Found {len(card_containers)} card containers")
            
            all_cards = []
            for i, container in enumerate(card_containers):
                logger.info(f"Processing card {i+1}/{len(card_containers)}")
                card_data = self.parse_card_data(container)
                if card_data:
                    all_cards.append(card_data)
                time.sleep(0.1)  # Be respectful to the server
                
            return all_cards
            
        except Exception as e:
            logger.error(f"Error scraping cards: {e}")
            return []
    
    def download_card_image(self, image_url: str, card_number: str, rarity: str) -> bool:
        """Download a card image"""
        try:
            # Create filename based on card number and rarity
            filename = f"{card_number}_{rarity}.jpg"
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
    
    def save_card_data(self, cards: List[Dict], filename: str = 'all_cards.json'):
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
        logger.info("Starting card scraping...")
        all_cards = self.scrape_all_cards()
        
        if not all_cards:
            logger.warning("No cards were scraped. Checking page structure...")
            # Let's try a more direct approach
            response = self.session.get(self.cardlist_url)
            soup = BeautifulSoup(response.content, 'html.parser')
            
            # Save the HTML for debugging
            with open('data/cardlist_page.html', 'w', encoding='utf-8') as f:
                f.write(soup.prettify())
            logger.info("Saved HTML for debugging")
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
                time.sleep(0.2)  # Be respectful
        
        logger.info("Scraping completed!")
        
        # Generate summary
        summary = {
            "total_cards": len(all_cards),
            "promo_cards": sum(1 for c in all_cards if c.get('is_promo', False)),
            "parallel_cards": sum(1 for c in all_cards if c.get('is_parallel', False)),
            "card_types": list(set(c.get('type', '') for c in all_cards if c.get('type'))),
            "colors": list(set(c.get('color', '') for c in all_cards if c.get('color'))),
            "rarities": list(set(c.get('rarity', '') for c in all_cards if c.get('rarity')))
        }
        
        with open('data/scraping_summary.json', 'w', encoding='utf-8') as f:
            json.dump(summary, f, ensure_ascii=False, indent=2)
        
        logger.info(f"Summary: {summary}")


if __name__ == "__main__":
    scraper = MeMeMeCardScraper()
    scraper.run()