<template>
  <div class="m-auto flex justify-between">
    <Card class="max-w-xl p-2">
      <CardHeader>
        <CardTitle>Product sale stats</CardTitle>
        <CardDescription>Show Top 5 Product sales</CardDescription>
      </CardHeader>
      <CardContent>
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Product Name</TableHead>
              <TableHead>Total Quantity Sold</TableHead>
              <TableHead>Total Amount</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            <TableRow v-if="isLoadingProducts">
              <TableCell colspan="3">
                <div class="h-60 w-full text-center">
                  ...Loading
                </div>
              </TableCell>
            </TableRow>
            <TableRow 
              v-else 
              v-for="(product, idx) in productSaleStats" 
              :key="idx" 
              data-testid="product-sales-rows"
            >
              <TableCell>{{ product.product_name  }}</TableCell>
              <TableCell class="text-center">{{ product.total_quantity_sold  }}</TableCell>
              <TableCell>${{ product.total_amount  }}</TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from '@/components/ui/card'

type TProductStats = {
  product_name: string
  total_quantity_sold: number
  total_amount: number
}
const isLoadingProducts = ref(false)
const productSaleStats = ref<TProductStats[]>()

onMounted(() => {
  fetchProductStats()
})

const fetchProductStats = async () => {
  try {
    isLoadingProducts.value = true
    const data = await fetch('http://localhost:8080/api/product-sale-stats').then(res => res.json()) as TProductStats[]
    productSaleStats.value = data
  } catch(e) {
    console.error(e)
  } finally {
    isLoadingProducts.value = false
  }
}
</script>
