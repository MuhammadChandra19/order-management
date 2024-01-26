<template>
  <div class="m-auto">
    <div class="flex gap-4 items-center my-2">
      <Search class="flex-none" width="16px"/>
      <div class="text-base flex-none">Search</div>
      <div class="flex-1 w-full">
        <Input 
          placeholder="Input by part of the order or product name"  
          v-model="searchText" 
          @keyup.enter="() => {
            $router.push({
              query: {
                ...queryParams,
                search: searchText
              }
            })
          }"
        />
      </div>

    </div>
    <DateRangePicker class="py-4" v-model="dateFilter"/>
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
                  @click="toggleSort"
                />
                <ArrowUpAz 
                  v-if="queryParams.sort_direction === 'ASC'" 
                  class="cursor-pointer" width="16px"
                  data-testid="sort-asc"
                  @click="toggleSort"
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
    <div class="flex w-full justify-center gap-4 items-center p-4">
      <div class="flex gap-2 items-center">
        <div class="text-sm">Total {{ orderList.length }}</div>
        <div class="p-2 border border-slate-600 rounded-sm" data-testid="dropdown-page-limit">
          <DropdownMenu>
            <DropdownMenuTrigger>{{ queryParams.limit }}/page</DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem @click="() => handleLimitChange(5)">5</DropdownMenuItem>
              <DropdownMenuItem @click="() => handleLimitChange(10)">10</DropdownMenuItem>
              <DropdownMenuItem @click="() => handleLimitChange(15)">15</DropdownMenuItem>
              <DropdownMenuItem @click="() => handleLimitChange(20)">20</DropdownMenuItem>
              <DropdownMenuItem @click="() => handleLimitChange(25)">25</DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
      <div class="flex gap-2 items-center">
        <ArrowLeftCircle width="16px" class="cursor-pointer" @click="() => handlePageChange(-1)"/>
        <div class="font-medium text-lg">
          {{ Math.ceil((queryParams.offset || 0) / (queryParams?.limit || 5)) + 1 }}
        </div>
        <ArrowRightCircle width="16px" class="cursor-pointer" @click="() => handlePageChange(1)"/>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { watch, onMounted, ref } from 'vue'
import format from 'date-fns-tz/format'
import { useRoute, useRouter } from 'vue-router';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import DateRangePicker, { TDateRangeValue } from '@/components/DateRangePicker.vue'
import Input from '@/components/ui/input/Input.vue'
import { ArrowDownZA, ArrowLeftCircle, ArrowRightCircle, ArrowUpAz, Search } from 'lucide-vue-next'
import {
  DropdownMenu,
  DropdownMenuItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'


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
const router = useRouter()

const orderList = ref<TOrderInfo[]>([])
const isLoading = ref<boolean>(false)
const queryParams = ref<{
  search?: string,
  start_date?: string,
  end_date?: string,
  sort_direction?: string,
  limit?: number,
  offset?: number
}>({ sort_direction: 'DESC', limit: 5, offset: 0 })

const dateFilter = ref<TDateRangeValue>()
const searchText = ref<string>()

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
    queryParams.value = route.query
    loadOrderList()
  }
)

watch(
  () => dateFilter,
  () => {
    if(dateFilter.value?.start && dateFilter.value.end) {
      router.push({
        query: {
          ...queryParams.value,
          start_date: `${formatDate(dateFilter.value?.start)}Z`,
          end_date: `${formatDate(dateFilter.value?.end)}Z`
        }
      })
    }
  },
  {
    deep: true
  }
)

const formatDate = (date: Date): string => {
  return format(date, 'yyyy-MM-dd hh:mm:ss', {
    timeZone: 'Australia/Melbourne'
  })
}

const toggleSort = () => {
  const isDescending = queryParams.value?.sort_direction === 'DESC' || false
  router.push({
    query: {
      ...queryParams.value,
      sort_direction: isDescending ? 'ASC' : 'DESC'
    }
  })
}

const handleLimitChange = (value: number) => {
  router.push({
    query: {
      ...queryParams.value,
      limit: value,
      offset: 0
    }
  })
}

const handlePageChange = (direction: number) => {
  const currentPage =  Math.ceil((queryParams.value.offset || 0) / (queryParams.value.limit || 5)) + direction
  const offset = currentPage > 0 ? currentPage * (queryParams.value.limit || 5) : 0
  router.push({
    query: {
      ...queryParams.value,
      offset
    }
  })
}

onMounted(() => {
  if (route.query?.start_date && route.query?.end_date) {
    dateFilter.value = {
      start: new Date(route.query.start_date as string),
      end: new Date(route.query.end_date as string)
    }
  }

  if(route.query.search) {
    searchText.value = route.query.search as string
  }
  loadOrderList()
})
  
</script>

