<script setup lang="ts">
import { useVModel } from '@vueuse/core'
import { format } from 'date-fns-tz'
import { Calendar as CalendarIcon } from 'lucide-vue-next'

import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { Calendar } from '@/components/ui/calendar'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'

export type TDateRangeValue = {
  start: Date,
  end: Date
}

const props = defineProps<{modelValue?: TDateRangeValue}>()
const emits = defineEmits<{
  (e: 'update:modelValue', payload: typeof props.modelValue): void
}>()

const modelValue = useVModel(props, 'modelValue', emits, {
  passive: true,
})

const formatDate = (date: Date): string => {
  return format(date, 'LLL dd, y', {
    timeZone: 'Australia/Melbourne'
  })
}
</script>

<template>
  <div :class="cn('grid gap-2', $attrs.class ?? '')">
    <Popover>
      <PopoverTrigger as-child>
        <Button
          id="date"
          :variant="'outline'"
          :class="cn(
            'w-[300px] justify-start text-left font-normal',
            !modelValue && 'text-muted-foreground',
          )"
        >
          <CalendarIcon class="mr-2 h-4 w-4" />

          <span>
            {{ modelValue?.start ? (
              modelValue.end ? `${formatDate(modelValue.start)} - ${formatDate(modelValue.end)}`
              : formatDate(modelValue.start)
            ) : 'Pick a date' }}
          </span>
        </Button>
      </PopoverTrigger>
      <PopoverContent class="w-auto p-0" align="start" :avoid-collisions="true">
        <Calendar
          v-model.range="modelValue"
          :columns="2"
        />
      </PopoverContent>
    </Popover>
  </div>
</template>
