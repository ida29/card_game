# MeMeMe TCG Card Data Download Report

## Summary
Successfully downloaded complete card data and images from https://mememe-tcg.com/cardlist

## Statistics
- **Total Cards Downloaded**: 124 cards
- **Promo Cards**: 12 cards (marked with (P) suffix)
- **Parallel Cards**: 10 cards (パラレル variants)
- **Card Images**: 119 unique images

## Card Types
- ふれんど (Friend)
- サポート (Support)
- フィールド (Field)

## Colors
- 赤 (Red)
- 青 (Blue)
- 黄 (Yellow)
- 緑 (Green)
- 無 (Colorless)

## Rarities
- C (Common)
- U (Uncommon)
- R (Rare)
- SR (Super Rare)
- SEC (Secret)
- Rパラレル (Rare Parallel)
- SRパラレル (Super Rare Parallel)
- SECパラレル (Secret Parallel)

## Data Fields Captured
Each card includes:
- Card number (including promo/parallel indicators)
- Card name
- Rarity
- Cost information:
  - Total cost
  - Color-specific costs (red, blue, yellow, green)
  - Colorless costs
- Card color
- Card type
- Attribute
- Emotion
- Power value
- Ability/effect text
- Flavor text (where available)
- Height/weight (for some cards)
- Image URL
- Promo/parallel flags

## File Locations
- **Card Data**: `/data/mememe_cards_complete.json`
- **Card Images**: `/data/card_images/`
- **Summary**: `/data/scraping_summary_final.json`

## Sample Card Data
```json
{
  "number": "F-001",
  "name": "バードン",
  "rarity": "C",
  "cost": {
    "total": 1,
    "red": 0,
    "blue": 0,
    "yellow": 0,
    "green": 0,
    "colorless": 1
  },
  "color": "赤",
  "type": "ふれんど",
  "attribute": "霊",
  "emotion": "喜",
  "power": 1000,
  "ability": "効果なし",
  "image_url": "https://mememe-tcg.com/assets/images/card/F-001_C.jpg"
}
```

## Notes
- All cards have complete cost breakdowns including total cost and color requirements
- Promo cards are identified by "(P)" in their number
- Parallel cards have "パラレル" in their rarity
- Images are saved with consistent naming: `{card_number}_{rarity}.jpg`