<script setup lang="ts">
import { watch, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import DateRangePicker from '@/components/DateRangePicker.vue'

import { ArrowDownZA } from 'lucide-vue-next'
type TOrderInfo = {
  order_name: string
  company_name: string
  customer_name: string
  product_name: string
  order_date: string
  delivered: string
  total_amount: number
}

const route = useRoute()

const orderList = ref<TOrderInfo[]>([])
const isLoading = ref<boolean>(false)

const loadOrderList = async () => {
  try {
    isLoading.value = true
    const data = await fetch(`http://localhost:8080/api${route.fullPath}`).then(res => res.json()) as TOrderInfo[]
    orderList.value = data
  } catch(e) {
    console.log(e)
  } finally {
    isLoading.value = false
  }
}

const parseDate = (dateStr: string) => {
  const date = new Date(dateStr);
  return new Intl.DateTimeFormat("en-US", {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  }).format(date);
}

watch(
  () => route.fullPath,
  () => {
    loadOrderList()
  },
)

onMounted(() => {
  loadOrderList()
})
  
</script>
<template>
  <div class="m-auto">
    <DateRangePicker class="py-4"/>
    <div class="py-2">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>
              Order name
            </TableHead>
            <TableHead>
              Customer Company
            </TableHead>
            <TableHead>
              Customer name
            </TableHead>
            <TableHead>
              <div class="flex gap-2">
                <div>Order date</div>
                <ArrowDownZA class="cursor-pointer" width="16px"/>
              </div>
            </TableHead>
            <TableHead>
              Delivered Amount
            </TableHead>
            <TableHead>
              Total Amount
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="isLoading">
            <TableCell colSpan="6">
              <div class="h-60 w-full text-center">
                ...Loading
              </div>
            </TableCell>
          </TableRow>
          <TableRow v-else v-for="(order, idx) in orderList" :key="idx">
            <TableCell class="font-medium">
              {{  order.order_name }}
            </TableCell>
            <TableCell>
              {{  order.company_name }}
            </TableCell>
            <TableCell>
              {{  order.customer_name }}
            </TableCell>
            <TableCell>
              {{  parseDate(order.order_date) }}
            </TableCell>
            <TableCell class="text-center">
              {{  order.delivered }}
            </TableCell>
            <TableCell class="text-center">
              ${{  order.total_amount }}
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  </div>
</template>

