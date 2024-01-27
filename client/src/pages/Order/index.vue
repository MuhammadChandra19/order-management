<template>
  <div class="m-auto">
    <div class="flex gap-4 items-center my-2">
      <Search class="flex-none" width="16px"/>
      <div class="text-base flex-none">Search</div>
      <div class="flex-1 w-full">
        <Input 
          data-testid="search-input"
          placeholder="Input by part of the order or product name"  
          v-model="searchText" 
          @keyup.enter="() => {
            if(searchText.trim().length > 0) {
              $router.push({
                query: {
                  ...queryParams,
                  search: searchText
                }
              })
            }
          }"
        />
      </div>

    </div>
    <DateRangePicker data-testid="date-range" class="py-4" v-model="dateFilter"/>
    <OrderTable 
      :is-loading="isLoading" 
      :query-params="queryParams" 
      :order-list="orderList" 
      @toggle-sort="toggleSort"
      @handle-page-change="handlePageChange"
      @handle-limit-change="handleLimitChange"
    />
  </div>
</template>
<script setup lang="ts">
import { watch, onMounted, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import DateRangePicker from '@/components/DateRangePicker.vue'
import Input from '@/components/ui/input/Input.vue'
import { Search } from 'lucide-vue-next'
import OrderTable from '@/components/OrderTable.vue';
import { TOrderFilter, TOrderInfo } from '@/lib/types';
import { formatDate } from '@/lib/utils';

const route = useRoute()
const router = useRouter()

const orderList = ref<TOrderInfo[]>([])
const isLoading = ref<boolean>(false)

const queryParams = computed(() => ({ 
  sort_direction: 'DESC', ...route.query 
}) as TOrderFilter)

const searchText = ref<string>(queryParams.value.search || '')

const dateFilter = computed({
  get() {
    const { value } = queryParams
    if (value.start_date && value.end_date) {
      return {
        start: new Date(value.start_date as string),
        end: new Date(value.end_date as string)
      }
    }
  }, 
  set(newValue) {
    if(newValue) {
      router.push({
        query: {
          ...queryParams.value,
          start_date: `${formatDate(newValue.start)}Z`,
          end_date: `${formatDate(newValue.end)}Z`
        }
      })
    }
  }
})


onMounted(() => {
  loadOrderList()
})

const loadOrderList = async () => {
  try {
    isLoading.value = true
    const data = await fetch(`http://localhost:8080/api${route.fullPath}`).then(res => res.json()) as TOrderInfo[]
    if(data) {
      orderList.value = data
    }
  } catch(e) {
    console.log(e)
  } finally {
    isLoading.value = false
  }
}

watch(
  () => route.fullPath,
  () => {
    loadOrderList()
  }
)

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
  
</script>

