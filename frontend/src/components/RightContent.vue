<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed} from 'vue'
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart,
  LineSeriesOption } from "echarts/charts";
import { Plus } from '@element-plus/icons-vue'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  ToolboxComponent,
  DataZoomComponent
} from "echarts/components";
import VChart from "vue-echarts";

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  GridComponent,
  TooltipComponent,
  LegendComponent,
  ToolboxComponent,
  DataZoomComponent
]);

const processingStatus = ref(2)
const activeStep = ref(0)

const countNumber = ref(3)
// const totalPercentage = ref(0)
const recordPercent = ref(0)
const prepareDataPercent = ref(0)
const anaysePercent = ref(0)
const percentDataHandle = ref(0)
const percentDataAnalyse = ref(0)
const responseTimeChartRef = ref()

const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]
const totalPercentage = computed(() => {
  return (recordPercent.value + prepareDataPercent.value + anaysePercent.value) / 3
}) 

const format = (percentage:number) => (percentage === 100 ? 'Done' : `${percentage}%`)

function startCountDown() {
  setInterval(() => {
    if (countNumber.value >0) {
      countNumber.value -=1
    }
  }, 1000)
}


const responseTimeData = reactive({
  title: {
    text: 'response time',
    left: "center",
  },
  tooltip: {
    trigger: 'axis',
  },
  grid: {
    left: '50px',
    right: '7px',
    bottom: '20px',
    top: '50px',
    containLabel: false
  },
  xAxis: {
    data: [],
    nameLocation: 'middle',
    max: function (value:any) {
      if (value.max < 10) {
        return 10
      }
      return value.max
    }
  },
  yAxis: {
    type: 'value',
    name: 'ms',
    position: 'left',
  },
  series: [
    {
      name: 'app',
      data: [],
      type: 'line',
      itemStyle: {
        color: 'rgb(46,211,111)'
      }
    }
  ]
})

function loadResponseTimeData(data: Array<number>) {
  responseTimeChartRef.value.setOption({ 
    series: [
      {
        data: data
      },
    ]
    })
}

function setRecordPercent(num: number){
  recordPercent.value = num
}

function setPrepareDataPercent(num: number){
  prepareDataPercent.value = num
}

function setAnaysePercent(num: number){
  anaysePercent.value = num
}

defineExpose({ 
  setRecordPercent,
  setPrepareDataPercent,
  setAnaysePercent,
  loadResponseTimeData 
})

</script>

<template>
  <div>
    <!-- <el-empty description="description" /> -->
    <div v-if="processingStatus===0">
       <span>
      操作说明
      - 选择设备
      - 打开待验证应用
      - 点击工具上准备按钮，此时倒计时3s后开始进入录制状态 （倒计时主要是考虑录制的启动时间及电脑操作后准备手机操作）
      - 录制10s 后自动结束并进行分析   （时长主要考虑1. 无需额外操作电脑， 2. 时间太长分析的时间也长）， 录制结束后关闭 调试信息
      - 分析后出多次操控结果，可考虑图表展示，excel 导出
    </span>
    </div>
     <div v-else-if="processingStatus===1">
      <span>
        倒计时：{{ countNumber }}
      </span>
     </div>
    <div v-else-if="processingStatus===2">
      <el-row class="progress-bar">
         <el-progress :percentage="totalPercentage" :format="format" />
      </el-row>

      <el-row>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="recordPercent" :color="colors">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">监听</span>
            </template>
          </el-progress>
        </el-col>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="prepareDataPercent" :color="colors">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">预处理</span>
            </template>
          </el-progress>
        </el-col>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="anaysePercent" :color="colors">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">分析</span>
            </template>
          </el-progress>
        </el-col>
      </el-row>
     
    </div>
    <div v-else-if="processingStatus===3">
      <span>
        处理完成
      </span>
      
    </div>
    <!-- <div style="width:500px;height:200px">
      <v-chart class="chart" ref="responseTimeChartRef" :option="responseTimeData"  />
    </div> -->
    <el-row style="width:(100vw-220px);height:200px">
       <v-chart class="chart" ref="responseTimeChartRef" :option="responseTimeData"  />
    </el-row>
    

  </div>

    
</template>

<style scoped>
.percentage-value {
  display: block;
  margin-top: 10px;
  font-size: 28px;
}
.percentage-label {
  display: block;
  margin-top: 10px;
  font-size: 12px;
}

.progress-bar .el-progress--line {
  margin-bottom: 15px;
  width: 100%;
}

</style>