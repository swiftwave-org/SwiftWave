<script setup>
import { formatTimestampHumannize } from '@/vendor/utils.js'

defineProps({
  title: {
    type: String,
    required: true
  },
  series: {
    type: Array,
    required: true
  },
  toolbar: {
    type: Boolean,
    required: false,
    default: true
  },
  yAxisFormatter: {
    type: Function,
    required: false,
    default: (value) => value
  },
  yAxisMinimumMax: {
    type: Number,
    default: 0,
    required: false
  }
})
</script>

<template>
  <div class="w-full">
    <p class="break-words text-sm font-medium">{{ title }}</p>
    <apexchart
      class="w-full"
      height="200"
      type="area"
      :options="{
        tooltip: {
          x: {
            formatter: formatTimestampHumannize
          }
        },
        animations: {
          enabled: false,
          initialAnimation: {
            enabled: false
          }
        },
        chart: {
          toolbar: {
            show: toolbar
          }
        },
        zoom: {
          autoScaleYaxis: false
        },
        markers: {
          size: 0
        },
        dataLabels: {
          enabled: false
        },
        yaxis: {
          tickAmount: 1,
          min: 0,
          max: (val) => Math.max(val, yAxisMinimumMax),
          labels: {
            show: true,
            formatter: yAxisFormatter
          }
        },
        xaxis: {
          type: 'datetime',
          tickAmount: 1,
          labels: {
            show: true
          }
        },
        stroke: {
          curve: 'smooth',
          width: 2
        },
        fill: {
          type: 'gradient',
          gradient: {
            shade: 'light',
            shadeIntensity: 1,
            opacityFrom: 0.2,
            opacityTo: 0.6,
            stops: [0, 100]
          }
        }
      }"
      :series="series"></apexchart>
  </div>
</template>

<style scoped></style>
