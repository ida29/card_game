#!/usr/bin/env python3
"""
Analyze discrepancies between existing card data and official card list
"""
import json
import csv
from collections import defaultdict

def load_existing_card_data():
    """Load existing card data from card_data.json"""
    with open('data/card_data.json', 'r', encoding='utf-8') as f:
        return json.load(f)

def load_official_card_data():
    """Load official card data"""
    with open('data/official_cardlist.json', 'r', encoding='utf-8') as f:
        return json.load(f)

def load_csv_data():
    """Load parsed CSV data"""
    cards = []
    with open('data/mememe_cards_parsed.csv', 'r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            cards.append(row)
    return cards

def normalize_card_number(number):
    """Normalize card number for comparison"""
    # Convert "F-001 (P)" to "F-001-P" format
    if " (P)" in number:
        return number.replace(" (P)", "-P")
    return number

def analyze_discrepancies():
    """Analyze discrepancies between data sources"""
    existing_data = load_existing_card_data()
    official_data = load_official_card_data()
    csv_data = load_csv_data()
    
    # Create lookup dictionaries
    existing_by_number = {normalize_card_number(card['number']): card for card in existing_data}
    official_by_number = {card['card_number']: card for card in official_data}
    
    # Track discrepancies
    discrepancies = {
        'missing_in_existing': [],
        'missing_in_official': [],
        'name_mismatch': [],
        'duplicate_entries': [],
        'promo_issues': []
    }
    
    # Check for duplicates in existing data
    number_count = defaultdict(int)
    for card in existing_data:
        normalized = normalize_card_number(card['number'])
        number_count[normalized] += 1
    
    for number, count in number_count.items():
        if count > 1:
            discrepancies['duplicate_entries'].append({
                'number': number,
                'count': count,
                'entries': [card for card in existing_data if normalize_card_number(card['number']) == number]
            })
    
    # Check cards in official data
    for card_number, official_card in official_by_number.items():
        if card_number not in existing_by_number:
            discrepancies['missing_in_existing'].append(official_card)
        else:
            existing_card = existing_by_number[card_number]
            # Check name consistency
            if existing_card['name'] != official_card['card_name']:
                discrepancies['name_mismatch'].append({
                    'number': card_number,
                    'existing_name': existing_card['name'],
                    'official_name': official_card['card_name']
                })
    
    # Check cards in existing data but not in official
    for card_number, existing_card in existing_by_number.items():
        if card_number not in official_by_number:
            # This might be because our official data is incomplete
            discrepancies['missing_in_official'].append(existing_card)
    
    # Analyze promo card issues
    for card in existing_data:
        if "-P" in card['number'] or " (P)" in card['number']:
            normalized = normalize_card_number(card['number'])
            # Check if promo has corresponding regular version
            base_number = normalized.replace("-P", "")
            if base_number not in existing_by_number and base_number not in official_by_number:
                discrepancies['promo_issues'].append({
                    'promo_number': card['number'],
                    'missing_base': base_number
                })
    
    return discrepancies

def generate_report(discrepancies):
    """Generate a detailed report of discrepancies"""
    report = []
    report.append("=== Card Data Discrepancy Analysis Report ===\n")
    
    # Duplicate entries
    if discrepancies['duplicate_entries']:
        report.append(f"\n## Duplicate Entries ({len(discrepancies['duplicate_entries'])} found):")
        for dup in discrepancies['duplicate_entries']:
            report.append(f"\n- Card {dup['number']} has {dup['count']} entries:")
            for entry in dup['entries']:
                report.append(f"  - Number: {entry['number']}, Name: {entry['name']}")
    
    # Missing in existing data
    if discrepancies['missing_in_existing']:
        report.append(f"\n## Missing in Existing Data ({len(discrepancies['missing_in_existing'])} cards):")
        for card in discrepancies['missing_in_existing']:
            report.append(f"- {card['card_number']}: {card['card_name']} ({card['rarity']}) - {card['card_type']}")
    
    # Name mismatches
    if discrepancies['name_mismatch']:
        report.append(f"\n## Name Mismatches ({len(discrepancies['name_mismatch'])} found):")
        for mismatch in discrepancies['name_mismatch']:
            report.append(f"- {mismatch['number']}: '{mismatch['existing_name']}' vs '{mismatch['official_name']}'")
    
    # Promo issues
    if discrepancies['promo_issues']:
        report.append(f"\n## Promo Card Issues ({len(discrepancies['promo_issues'])} found):")
        for issue in discrepancies['promo_issues']:
            report.append(f"- Promo {issue['promo_number']} missing base card {issue['missing_base']}")
    
    # Summary
    report.append("\n## Summary:")
    report.append(f"- Total duplicate entries: {len(discrepancies['duplicate_entries'])}")
    report.append(f"- Total missing cards: {len(discrepancies['missing_in_existing'])}")
    report.append(f"- Total name mismatches: {len(discrepancies['name_mismatch'])}")
    report.append(f"- Total promo issues: {len(discrepancies['promo_issues'])}")
    
    return "\n".join(report)

def main():
    """Main function"""
    print("Analyzing card data discrepancies...")
    
    try:
        discrepancies = analyze_discrepancies()
        report = generate_report(discrepancies)
        
        # Save report to file
        with open('data/discrepancy_report.txt', 'w', encoding='utf-8') as f:
            f.write(report)
        
        # Print report
        print(report)
        
        # Save detailed JSON for further analysis
        with open('data/discrepancies.json', 'w', encoding='utf-8') as f:
            json.dump(discrepancies, f, ensure_ascii=False, indent=2)
        
        print("\nDetailed analysis saved to:")
        print("- data/discrepancy_report.txt")
        print("- data/discrepancies.json")
        
    except Exception as e:
        print(f"Error: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    main()