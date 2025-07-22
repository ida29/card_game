#!/usr/bin/env python3
"""
Comprehensive analysis of card data discrepancies
"""
import json
import re
from collections import defaultdict

# Official card list from website (complete listing)
OFFICIAL_CARDS = {
    # Red Friends (赤 ふれんど)
    "F-001": {"name": "バードン", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 1, "power": 1000},
    "F-002": {"name": "なみだぶくろん", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 2, "power": 2000},
    "F-003": {"name": "フラフラ", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 2, "power": 2000},
    "F-004": {"name": "ハシルシト", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 2, "power": 2000},
    "F-005": {"name": "かうちゃん", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 3, "power": 3000},
    "F-006": {"name": "ヒヤケラトプス", "rarity": "U", "type": "ふれんど", "color": "赤", "cost": 3, "power": 4000},
    "F-007": {"name": "るくそー", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 3, "power": 3000},
    "F-008": {"name": "ボーイ", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 3, "power": 3000},
    "F-009": {"name": "八つ目", "rarity": "U", "type": "ふれんど", "color": "赤", "cost": 4, "power": 5000},
    "F-010": {"name": "ブロントくん", "rarity": "U", "type": "ふれんど", "color": "赤", "cost": 4, "power": 5000},
    "F-011": {"name": "ポチ", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 4, "power": 4000},
    "F-012": {"name": "くらげ坊", "rarity": "C", "type": "ふれんど", "color": "赤", "cost": 4, "power": 4000},
    "F-013": {"name": "るくそー", "rarity": "R", "type": "ふれんど", "color": "赤", "cost": 4, "power": 5000},
    "F-013-P": {"name": "るくそー", "rarity": "R", "type": "ふれんど", "color": "赤", "cost": 4, "power": 5000},
    "F-014": {"name": "くらげ坊", "rarity": "U", "type": "ふれんど", "color": "赤", "cost": 5, "power": 6000},
    "F-015": {"name": "ティラノちゃん", "rarity": "R", "type": "ふれんど", "color": "赤", "cost": 5, "power": 7000},
    "F-015-P": {"name": "ティラノちゃん", "rarity": "R", "type": "ふれんど", "color": "赤", "cost": 5, "power": 7000},
    "F-016": {"name": "くらげ坊", "rarity": "SR", "type": "ふれんど", "color": "赤", "cost": 8, "power": 12000},
    "F-016-P": {"name": "くらげ坊", "rarity": "SR", "type": "ふれんど", "color": "赤", "cost": 8, "power": 12000},
    
    # Blue Friends (青 ふれんど)
    "F-017": {"name": "じゅごん", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 1, "power": 1000},
    "F-018": {"name": "イソチャック", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 2, "power": 2000},
    "F-019": {"name": "カラカッサン", "rarity": "U", "type": "ふれんど", "color": "青", "cost": 2, "power": 3000},
    "F-020": {"name": "マルカニ", "rarity": "U", "type": "ふれんど", "color": "青", "cost": 2, "power": 3000},
    "F-021": {"name": "イチゴバット", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 2, "power": 2000},
    "F-022": {"name": "ジョニー", "rarity": "R", "type": "ふれんど", "color": "青", "cost": 2, "power": 3000},
    "F-023": {"name": "ユピ", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 1, "power": 1000},
    "F-023-P": {"name": "ユピ", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 1, "power": 1000},
    "F-024": {"name": "ジェラチンポン", "rarity": "U", "type": "ふれんど", "color": "青", "cost": 3, "power": 4000},
    "F-025": {"name": "しもん", "rarity": "R", "type": "ふれんど", "color": "青", "cost": 3, "power": 4000},
    "F-025-P": {"name": "しもん", "rarity": "R", "type": "ふれんど", "color": "青", "cost": 3, "power": 4000},
    "F-026": {"name": "ランプ", "rarity": "U", "type": "ふれんど", "color": "青", "cost": 3, "power": 4000},
    "F-027": {"name": "ハンマーヘッドくん", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 3, "power": 3000},
    "F-028": {"name": "ビッグアーム", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 4, "power": 4000},
    "F-029": {"name": "カラカッサン", "rarity": "R", "type": "ふれんど", "color": "青", "cost": 4, "power": 5000},
    "F-030": {"name": "クッキー", "rarity": "C", "type": "ふれんど", "color": "青", "cost": 4, "power": 4000},
    "F-031": {"name": "ユピ", "rarity": "R", "type": "ふれんど", "color": "青", "cost": 5, "power": 7000},
    "F-032": {"name": "しもん", "rarity": "SR", "type": "ふれんど", "color": "青", "cost": 6, "power": 9000},
    "F-032-P": {"name": "しもん", "rarity": "SR", "type": "ふれんど", "color": "青", "cost": 6, "power": 9000},
    
    # Yellow Friends (黄 ふれんど)
    "F-033": {"name": "ボーグ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 1000},
    "F-034": {"name": "メガロッコ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 500},
    "F-034-P": {"name": "メガロッコ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 500},
    "F-035": {"name": "ジセダイ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 1000},
    "F-036": {"name": "わらべ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 1000},
    "F-037": {"name": "カンブリアン", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 2, "power": 2000},
    "F-038": {"name": "シンゴーくん", "rarity": "U", "type": "ふれんど", "color": "黄", "cost": 2, "power": 3000},
    "F-039": {"name": "ヌリカベーゼ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 2, "power": 2000},
    "F-040": {"name": "フジさん", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 3, "power": 3000},
    "F-041": {"name": "ハヤオ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 1000},
    "F-041-P": {"name": "ハヤオ", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 1, "power": 1000},
    "F-042": {"name": "うっきー", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 2, "power": 2000},
    "F-042-P": {"name": "うっきー", "rarity": "C", "type": "ふれんど", "color": "黄", "cost": 2, "power": 2000},
    "F-043": {"name": "メガロッコ", "rarity": "R", "type": "ふれんど", "color": "黄", "cost": 3, "power": 4000},
    "F-044": {"name": "うっきー", "rarity": "R", "type": "ふれんど", "color": "黄", "cost": 4, "power": 5000},
    "F-045": {"name": "FM94.3", "rarity": "U", "type": "ふれんど", "color": "黄", "cost": 3, "power": 4000},
    "F-046": {"name": "ピンヒール", "rarity": "U", "type": "ふれんど", "color": "黄", "cost": 4, "power": 5000},
    "F-047": {"name": "だんぼー", "rarity": "U", "type": "ふれんど", "color": "黄", "cost": 4, "power": 5000},
    "F-048": {"name": "ハヤオ", "rarity": "SR", "type": "ふれんど", "color": "黄", "cost": 5, "power": 8000},
    "F-048-P": {"name": "ハヤオ", "rarity": "SR", "type": "ふれんど", "color": "黄", "cost": 5, "power": 8000},
    
    # Green Friends (緑 ふれんど)
    "F-049": {"name": "くるくるるん", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 1, "power": 1000},
    "F-050": {"name": "みのたま", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 1, "power": 1000},
    "F-051": {"name": "クワクワクワ", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 2, "power": 2000},
    "F-052": {"name": "ゆーゆーれん", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 2, "power": 2000},
    "F-053": {"name": "イガイガリン", "rarity": "U", "type": "ふれんど", "color": "緑", "cost": 2, "power": 3000},
    "F-054": {"name": "オオケムシ", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 2, "power": 2000},
    "F-055": {"name": "Ko2", "rarity": "R", "type": "ふれんど", "color": "緑", "cost": 3, "power": 3000},
    "F-055-P": {"name": "Ko2", "rarity": "R", "type": "ふれんど", "color": "緑", "cost": 3, "power": 3000},
    "F-056": {"name": "シーラン", "rarity": "U", "type": "ふれんど", "color": "緑", "cost": 2, "power": 3000},
    "F-056-P": {"name": "シーラン", "rarity": "U", "type": "ふれんど", "color": "緑", "cost": 2, "power": 3000},
    "F-057": {"name": "モモイ", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 3, "power": 3000},
    "F-058": {"name": "シーラン", "rarity": "R", "type": "ふれんど", "color": "緑", "cost": 4, "power": 5000},
    "F-059": {"name": "オニガワラン", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 4, "power": 4000},
    "F-060": {"name": "カゲロウくん", "rarity": "U", "type": "ふれんど", "color": "緑", "cost": 4, "power": 5000},
    "F-061": {"name": "ペチカ", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 4, "power": 4000},
    "F-062": {"name": "えびすけ", "rarity": "C", "type": "ふれんど", "color": "緑", "cost": 5, "power": 5000},
    "F-063": {"name": "トーテムポーラ", "rarity": "U", "type": "ふれんど", "color": "緑", "cost": 5, "power": 6000},
    "F-064": {"name": "Ko2", "rarity": "SR", "type": "ふれんど", "color": "緑", "cost": 7, "power": 10000},
    "F-064-P": {"name": "Ko2", "rarity": "SR", "type": "ふれんど", "color": "緑", "cost": 7, "power": 10000},
    
    # Support Cards (サポート)
    "F-065": {"name": "ガーディアン", "rarity": "C", "type": "サポート", "color": "赤", "cost": 1},
    "F-066": {"name": "正志とくらげ坊", "rarity": "C", "type": "サポート", "color": "赤", "cost": 1},
    "F-066-P": {"name": "正志とくらげ坊", "rarity": "C", "type": "サポート", "color": "赤", "cost": 1},
    "F-067": {"name": "大好物！", "rarity": "C", "type": "サポート", "color": "赤", "cost": 2},
    "F-068": {"name": "デコーレーション", "rarity": "C", "type": "サポート", "color": "赤", "cost": 2},
    "F-069": {"name": "特異点が開く扉", "rarity": "U", "type": "サポート", "color": "赤", "cost": 3},
    "F-070": {"name": "ブルードラゴン飛連蹴", "rarity": "R", "type": "サポート", "color": "赤", "cost": 4},
    "F-070-P": {"name": "ブルードラゴン飛連蹴", "rarity": "R", "type": "サポート", "color": "赤", "cost": 4},
    "F-071": {"name": "絶対に裏切らない友達", "rarity": "C", "type": "サポート", "color": "青", "cost": 1},
    "F-072": {"name": "竜也とユピ", "rarity": "C", "type": "サポート", "color": "青", "cost": 1},
    "F-073": {"name": "古池ダイビング", "rarity": "C", "type": "サポート", "color": "青", "cost": 2},
    "F-074": {"name": "出現", "rarity": "U", "type": "サポート", "color": "青", "cost": 2},
    "F-075": {"name": "同調攻撃！", "rarity": "C", "type": "サポート", "color": "青", "cost": 3},
    "F-076": {"name": "ヘッドバット", "rarity": "R", "type": "サポート", "color": "青", "cost": 3},
    "F-076-P": {"name": "ヘッドバット", "rarity": "R", "type": "サポート", "color": "青", "cost": 3},
    "F-077": {"name": "チーム対抗", "rarity": "C", "type": "サポート", "color": "黄", "cost": 1},
    "F-078": {"name": "念力", "rarity": "C", "type": "サポート", "color": "黄", "cost": 1},
    "F-079": {"name": "聞いちゃった！", "rarity": "C", "type": "サポート", "color": "黄", "cost": 2},
    "F-080": {"name": "謎の四人衆", "rarity": "C", "type": "サポート", "color": "黄", "cost": 3},
    "F-080-P": {"name": "謎の四人衆", "rarity": "C", "type": "サポート", "color": "黄", "cost": 3},
    "F-081": {"name": "キャベン・ディッシュ", "rarity": "U", "type": "サポート", "color": "黄", "cost": 3},
    "F-082": {"name": "絶対否定領域", "rarity": "R", "type": "サポート", "color": "黄", "cost": 4},
    "F-082-P": {"name": "絶対否定領域", "rarity": "R", "type": "サポート", "color": "黄", "cost": 4},
    "F-083": {"name": "巧とゲーム", "rarity": "C", "type": "サポート", "color": "緑", "cost": 1},
    "F-084": {"name": "ワイルドワインド", "rarity": "C", "type": "サポート", "color": "緑", "cost": 1},
    "F-085": {"name": "プログラム解析", "rarity": "C", "type": "サポート", "color": "緑", "cost": 2},
    "F-086": {"name": "最強ふれんど決定戦", "rarity": "C", "type": "サポート", "color": "緑", "cost": 2},
    "F-087": {"name": "熱風破砕掌", "rarity": "R", "type": "サポート", "color": "緑", "cost": 3},
    "F-087-P": {"name": "熱風破砕掌", "rarity": "R", "type": "サポート", "color": "緑", "cost": 3},
    "F-088": {"name": "「ゲーム終了」", "rarity": "U", "type": "サポート", "color": "緑", "cost": 4},
    
    # Field Cards (フィールド)
    "F-089": {"name": "見晴らし台", "rarity": "C", "type": "フィールド", "color": "赤", "cost": 1},
    "F-090": {"name": "神社", "rarity": "C", "type": "フィールド", "color": "赤", "cost": 1},
    "F-091": {"name": "正志の家", "rarity": "U", "type": "フィールド", "color": "赤", "cost": 1},
    "F-092": {"name": "不思議な教室", "rarity": "C", "type": "フィールド", "color": "青", "cost": 1},
    "F-093": {"name": "学園のプール", "rarity": "C", "type": "フィールド", "color": "青", "cost": 1},
    "F-094": {"name": "研究所", "rarity": "U", "type": "フィールド", "color": "青", "cost": 1},
    "F-095": {"name": "都雲大学", "rarity": "C", "type": "フィールド", "color": "黄", "cost": 1},
    "F-096": {"name": "天地救世教会本部", "rarity": "C", "type": "フィールド", "color": "黄", "cost": 1},
    "F-097": {"name": "都雲祭", "rarity": "U", "type": "フィールド", "color": "黄", "cost": 1},
    "F-098": {"name": "もぐらの家", "rarity": "C", "type": "フィールド", "color": "緑", "cost": 1},
    "F-099": {"name": "ごみ捨て場", "rarity": "C", "type": "フィールド", "color": "緑", "cost": 1},
    "F-100": {"name": "屋上", "rarity": "U", "type": "フィールド", "color": "緑", "cost": 1},
    
    # Secret Cards
    "F-101": {"name": "オーバル", "rarity": "SEC", "type": "ふれんど", "color": "無", "cost": 10, "power": 15000},
    "F-101-P": {"name": "オーバル", "rarity": "SEC", "type": "ふれんど", "color": "無", "cost": 10, "power": 15000},
    "F-102": {"name": "くらげ坊(変身)", "rarity": "SEC", "type": "ふれんど", "color": "赤", "cost": 7, "power": 11000},
    "F-102-P": {"name": "くらげ坊(変身)", "rarity": "SEC", "type": "ふれんど", "color": "赤", "cost": 7, "power": 11000},
}

def normalize_card_number(number):
    """Normalize card number for comparison"""
    if " (P)" in number:
        return number.replace(" (P)", "-P")
    return number

def load_existing_data():
    """Load existing card data"""
    with open('data/card_data.json', 'r', encoding='utf-8') as f:
        return json.load(f)

def analyze_all_discrepancies():
    """Comprehensive analysis of all discrepancies"""
    existing_data = load_existing_data()
    
    # Create lookup dictionary for existing data
    existing_by_number = {}
    for card in existing_data:
        normalized = normalize_card_number(card['number'])
        if normalized not in existing_by_number:
            existing_by_number[normalized] = []
        existing_by_number[normalized].append(card)
    
    discrepancies = {
        'missing_cards': [],
        'duplicate_entries': [],
        'name_mismatches': [],
        'wrong_rarity': [],
        'promo_format_issues': [],
        'image_issues': []
    }
    
    # Check for missing cards
    for card_number, official_card in OFFICIAL_CARDS.items():
        if card_number not in existing_by_number:
            discrepancies['missing_cards'].append({
                'number': card_number,
                'name': official_card['name'],
                'type': official_card['type'],
                'rarity': official_card['rarity']
            })
    
    # Check existing cards
    for normalized_number, cards in existing_by_number.items():
        # Check for duplicates
        if len(cards) > 1:
            discrepancies['duplicate_entries'].append({
                'number': normalized_number,
                'count': len(cards),
                'entries': [{'number': c['number'], 'name': c['name']} for c in cards]
            })
        
        # Check against official data
        if normalized_number in OFFICIAL_CARDS:
            official = OFFICIAL_CARDS[normalized_number]
            for card in cards:
                # Check name
                if card['name'] != official['name']:
                    discrepancies['name_mismatches'].append({
                        'number': normalized_number,
                        'existing_name': card['name'],
                        'official_name': official['name']
                    })
                
                # Check image naming convention
                expected_rarity = official['rarity']
                image_path = card.get('local_image_path', '')
                
                # Check if image filename matches expected pattern
                if normalized_number.endswith('-P'):
                    # Promo cards should have special naming
                    if '-P_' not in image_path:
                        discrepancies['image_issues'].append({
                            'number': normalized_number,
                            'issue': 'Promo card image naming incorrect',
                            'current': image_path,
                            'expected_pattern': f"{normalized_number}_{expected_rarity}_Color.jpg"
                        })
                else:
                    # Regular cards
                    expected_pattern = f"{normalized_number}_{expected_rarity}.jpg"
                    if expected_pattern not in image_path:
                        discrepancies['image_issues'].append({
                            'number': normalized_number,
                            'issue': 'Image filename doesn\'t match expected pattern',
                            'current': image_path,
                            'expected': expected_pattern
                        })
        
        # Check promo format issues
        for card in cards:
            if "(P)" in card['number'] or "-P" in card['number']:
                if "(P)" in card['number'] and "-P" in card['number']:
                    discrepancies['promo_format_issues'].append({
                        'number': card['number'],
                        'issue': 'Inconsistent promo format'
                    })
    
    return discrepancies

def generate_detailed_report(discrepancies):
    """Generate detailed report with recommendations"""
    report = []
    report.append("=== COMPREHENSIVE CARD DATA ANALYSIS REPORT ===")
    report.append(f"\nAnalysis Date: 2025-07-22")
    report.append(f"Total Official Cards: {len(OFFICIAL_CARDS)}")
    
    # Missing Cards
    if discrepancies['missing_cards']:
        report.append(f"\n## MISSING CARDS ({len(discrepancies['missing_cards'])} cards)")
        report.append("These cards exist in the official list but are missing from our database:")
        
        # Group by type
        by_type = defaultdict(list)
        for card in discrepancies['missing_cards']:
            by_type[card['type']].append(card)
        
        for card_type, cards in by_type.items():
            report.append(f"\n### {card_type} ({len(cards)} cards):")
            for card in sorted(cards, key=lambda x: x['number']):
                report.append(f"  - {card['number']}: {card['name']} ({card['rarity']})")
    
    # Duplicate Entries
    if discrepancies['duplicate_entries']:
        report.append(f"\n## DUPLICATE ENTRIES ({len(discrepancies['duplicate_entries'])} found)")
        report.append("These card numbers have multiple entries in the database:")
        for dup in discrepancies['duplicate_entries']:
            report.append(f"\n- {dup['number']} appears {dup['count']} times:")
            for entry in dup['entries']:
                report.append(f"  - Number: {entry['number']}, Name: {entry['name']}")
    
    # Name Mismatches
    if discrepancies['name_mismatches']:
        report.append(f"\n## NAME MISMATCHES ({len(discrepancies['name_mismatches'])} found)")
        for mismatch in discrepancies['name_mismatches']:
            report.append(f"- {mismatch['number']}: '{mismatch['existing_name']}' → '{mismatch['official_name']}'")
    
    # Image Issues
    if discrepancies['image_issues']:
        report.append(f"\n## IMAGE NAMING ISSUES ({len(discrepancies['image_issues'])} found)")
        for issue in discrepancies['image_issues'][:10]:  # Show first 10
            report.append(f"- {issue['number']}: {issue['issue']}")
            report.append(f"  Current: {issue['current']}")
            if 'expected' in issue:
                report.append(f"  Expected: {issue['expected']}")
    
    # Summary and Recommendations
    report.append("\n## SUMMARY")
    report.append(f"- Missing cards: {len(discrepancies['missing_cards'])}")
    report.append(f"- Duplicate entries: {len(discrepancies['duplicate_entries'])}")
    report.append(f"- Name mismatches: {len(discrepancies['name_mismatches'])}")
    report.append(f"- Image naming issues: {len(discrepancies['image_issues'])}")
    
    report.append("\n## RECOMMENDATIONS")
    report.append("1. Add all missing cards to the database")
    report.append("2. Remove duplicate entries (keep only one per card number)")
    report.append("3. Standardize promo card format to use -P suffix")
    report.append("4. Fix image naming to match pattern: {number}_{rarity}.jpg")
    report.append("5. Update card names to match official names")
    
    return "\n".join(report)

def main():
    """Main analysis function"""
    print("Running comprehensive card data analysis...")
    
    discrepancies = analyze_all_discrepancies()
    report = generate_detailed_report(discrepancies)
    
    # Save report
    with open('data/comprehensive_analysis_report.txt', 'w', encoding='utf-8') as f:
        f.write(report)
    
    # Save detailed JSON
    with open('data/comprehensive_discrepancies.json', 'w', encoding='utf-8') as f:
        json.dump(discrepancies, f, ensure_ascii=False, indent=2)
    
    print(report)
    print("\n\nAnalysis complete! Files saved:")
    print("- data/comprehensive_analysis_report.txt")
    print("- data/comprehensive_discrepancies.json")

if __name__ == "__main__":
    main()