export type TOrderInfo = {
  order_name: string
  company_name: string
  customer_name: string
  product_name: string
  order_date: string
  delivered: string
  total_amount: number
}

export type TOrderFilter = {
  search?: string,
  start_date?: string,
  end_date?: string,
  sort_direction?: string,
  limit?: number,
  offset?: number
}