<template>
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
              <ArrowDownZA
                v-if="queryParams.sort_direction === 'DESC'" 
                class="cursor-pointer" width="16px"
                data-testid="sort-desc"
                @click="$emit('toggleSort')"
              />
              <ArrowUpAZ
                v-if="queryParams.sort_direction === 'ASC'" 
                class="cursor-pointer" width="16px"
                data-testid="sort-asc"
                @click="$emit('toggleSort')"
              />
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
        <TableRow v-else-if="!isLoading && orderList.length > 0" v-for="(order, idx) in orderList" :key="idx">
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
        <TableRow v-else>
          <TableCell colSpan="6" class="h-60 w-full text-center">No data</TableCell>
        </TableRow>
      </TableBody>
    </Table>
  </div>
  <div class="flex w-full justify-center gap-4 items-center p-4">
    <div class="flex gap-2 items-center">
      <div class="text-sm">Total {{ orderList.length }}</div>
      <div class="p-2 border border-slate-600 rounded-sm" data-testid="dropdown-page-limit">
        <DropdownMenu>
          <DropdownMenuTrigger >{{ queryParams.limit || 5 }}/page</DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuItem @click="() => $emit('handleLimitChange',5)">5</DropdownMenuItem>
            <DropdownMenuItem @click="() => $emit('handleLimitChange',10)">10</DropdownMenuItem>
            <DropdownMenuItem @click="() => $emit('handleLimitChange',15)">15</DropdownMenuItem>
            <DropdownMenuItem @click="() => $emit('handleLimitChange',20)">20</DropdownMenuItem>
            <DropdownMenuItem @click="() => $emit('handleLimitChange',25)">25</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
    <div class="flex gap-2 items-center">
      <ArrowLeftCircle 
        width="16px" 
        class="cursor-pointer" 
        data-testid="previous-page-arrow" 
        @click="() => $emit('handlePageChange', -1)"
      />
      <div class="font-medium text-lg">
        {{ Math.ceil((queryParams.offset || 0) / (queryParams?.limit || 5)) + 1 }}
      </div>
      <ArrowRightCircle
        width="16px" 
        data-testid="next-page-arrow" 
        class="cursor-pointer"
        @click="() => $emit('handlePageChange', 1)"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from '@/components/ui/table'
import { TOrderFilter, TOrderInfo } from '@/lib/types';
import {
  DropdownMenu,
  DropdownMenuItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import { parseDate } from '@/lib/utils'
import { ArrowDownZA, ArrowLeftCircle, ArrowRightCircle, ArrowUpAZ } from 'lucide-vue-next';

defineProps<{
  queryParams: TOrderFilter,
  isLoading: boolean,
  orderList: TOrderInfo[]
}>()

defineEmits<{
  (e: 'toggleSort'): void,
  (e: 'handlePageChange', direction: number): void,
  (e: 'handleLimitChange', limit: number): void
}>()
</script>