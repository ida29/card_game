#!/usr/bin/env python3
"""
Fix identified card data discrepancies
"""
import json
import os
import shutil

def load_card_data():
    """Load existing card data"""
    with open('data/card_data.json', 'r', encoding='utf-8') as f:
        return json.load(f)

def save_card_data(data):
    """Save updated card data"""
    # Backup existing data
    shutil.copy('data/card_data.json', 'data/card_data_backup_fix.json')
    
    with open('data/card_data.json', 'w', encoding='utf-8') as f:
        json.dump(data, f, ensure_ascii=False, indent=4)

def fix_duplicate_entries(cards):
    """Remove duplicate F-016-P entry"""
    print("Fixing duplicate entries...")
    
    # Keep track of seen numbers
    seen = set()
    cleaned_cards = []
    
    for card in cards:
        # Normalize number for checking
        normalized = card['number'].replace(' (P)', '-P')
        
        # Skip if we've already seen this card
        if normalized in seen:
            print(f"  Removing duplicate: {card['number']} - {card['name']}")
            continue
        
        seen.add(normalized)
        
        # Standardize promo format
        if ' (P)' in card['number']:
            card['number'] = card['number'].replace(' (P)', '-P')
            print(f"  Standardized promo format: {card['name']} to {card['number']}")
        
        cleaned_cards.append(card)
    
    return cleaned_cards

def fix_name_mismatches(cards):
    """Fix card name mismatches"""
    print("\nFixing name mismatches...")
    
    name_fixes = {
        "F-068": "デコーレーション"  # Fix long dash to regular dash
    }
    
    for card in cards:
        number = card['number'].replace('-P', '').replace(' (P)', '')
        if number in name_fixes and card['name'] != name_fixes[number]:
            old_name = card['name']
            card['name'] = name_fixes[number]
            print(f"  Fixed name for {card['number']}: '{old_name}' → '{card['name']}'")
    
    return cards

def standardize_image_paths(cards):
    """Standardize image paths for consistency"""
    print("\nStandardizing image paths...")
    
    # Map of correct image naming patterns
    image_fixes = {
        "F-032-P": "F-032_SR-P.jpg",
        "F-048-P": "F-048_SR-P.jpg", 
        "F-064-P": "F-064_SR_P.jpg",
        "F-070-P": "F-070_R-P.jpg",
        "F-076-P": "F-076_R-P.jpg",
        "F-082-P": "F-082_R-P.jpg",
        "F-087-P": "F-087_R-P.jpg",
        "F-101-P": "F-101_SEC_P.jpg",
        "F-102-P": "F-102_SEC_P.jpg"
    }
    
    for card in cards:
        normalized_number = card['number'].replace(' (P)', '-P')
        
        # Update image URLs for promo cards
        if normalized_number in image_fixes:
            expected_filename = image_fixes[normalized_number]
            card['image_url'] = f"https://mememe-tcg.com/assets/images/card/{expected_filename}"
            card['local_image_path'] = f"card_images/{expected_filename}"
            print(f"  Updated image paths for {normalized_number}")
    
    return cards

def verify_all_cards_present(cards):
    """Verify we have all expected cards"""
    print("\nVerifying card completeness...")
    
    expected_cards = set()
    # Regular cards F-001 to F-102
    for i in range(1, 103):
        expected_cards.add(f"F-{i:03d}")
    
    # Known promo cards
    promo_cards = [
        "F-013-P", "F-015-P", "F-016-P", "F-023-P", "F-025-P",
        "F-032-P", "F-034-P", "F-041-P", "F-042-P", "F-048-P",
        "F-055-P", "F-056-P", "F-064-P", "F-066-P", "F-070-P",
        "F-076-P", "F-080-P", "F-082-P", "F-087-P", "F-101-P", "F-102-P"
    ]
    expected_cards.update(promo_cards)
    
    # Check what we have
    existing_numbers = {card['number'].replace(' (P)', '-P') for card in cards}
    
    missing = expected_cards - existing_numbers
    if missing:
        print(f"  WARNING: Missing {len(missing)} cards: {sorted(missing)}")
    else:
        print(f"  ✓ All {len(expected_cards)} expected cards are present")
    
    extra = existing_numbers - expected_cards
    if extra:
        print(f"  INFO: Extra cards found: {sorted(extra)}")

def main():
    """Main function to fix all discrepancies"""
    print("Loading card data...")
    cards = load_card_data()
    print(f"Loaded {len(cards)} cards")
    
    # Apply fixes
    cards = fix_duplicate_entries(cards)
    cards = fix_name_mismatches(cards)
    cards = standardize_image_paths(cards)
    
    # Sort cards by number for consistency
    def sort_key(card):
        number = card['number'].replace('F-', '').replace('-P', '.5')
        try:
            return float(number)
        except:
            return 999
    
    cards.sort(key=sort_key)
    
    # Verify completeness
    verify_all_cards_present(cards)
    
    # Save updated data
    print(f"\nSaving {len(cards)} cards...")
    save_card_data(cards)
    
    print("\n✓ All fixes applied successfully!")
    print("  Backup saved to: data/card_data_backup_fix.json")
    print("  Updated data saved to: data/card_data.json")

if __name__ == "__main__":
    main()