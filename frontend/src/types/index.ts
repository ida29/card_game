export type CardType = 'ふれんど' | 'サポート' | 'フィールド'
export type CardColor = '赤' | '青' | '黄' | '緑'
export type CardRarity = 'C' | 'U' | 'R' | 'SR' | 'SEC' | 'C-P' | 'U-P' | 'R-P' | 'SR-P' | 'SEC-P'

export interface Card {
  ID: number
  card_no: string
  name: string
  type: CardType
  color: CardColor
  cost: number
  cost_red: number
  cost_blue: number
  cost_yellow: number
  cost_green: number
  cost_colorless: number
  power?: number
  rarity: CardRarity
  effect?: string
  flavor_text?: string
  image_url: string
  local_image_path: string
  energy_icons?: string[]
  is_counter: boolean
  is_main_counter: boolean
  is_promo?: boolean
}

export interface DeckCard {
  ID: number
  deck_id: number
  card_no: string
  quantity: number
  card: Card
}

export interface Deck {
  ID: number
  name: string
  user_id: number
  cards: DeckCard[]
  is_active: boolean
}

export interface GameState {
  player1_state: PlayerState
  player2_state: PlayerState
}

export interface PlayerState {
  deck: string[]
  hand: string[]
  battle_area: Record<string, Friend>
  energy_area: EnergyCard[]
  negative_energy: string[]
  trash: string[]
  field_card?: string
}

export interface Friend {
  card_no: string
  power: number
  is_rest: boolean
  turn_played: number
}

export interface EnergyCard {
  card_no: string
  color: CardColor
  is_rest: boolean
}