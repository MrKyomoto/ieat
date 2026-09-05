export interface FoodWindow {
  id: string
  externalCode: string
  name: string
  description: string
  businessHours: string
}

export interface Floor {
  id: string
  name: string
  windows: FoodWindow[]
}

export interface Canteen {
  id: string
  name: string
  floors: Floor[]
}
