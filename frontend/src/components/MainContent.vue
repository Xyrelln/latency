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

const processingStatus = ref(3)

const countNumber = ref(3)
const percentRecord = ref(60)
const percentDataHandle = ref(60)
const percentDataAnalyse = ref(20)
const responseTimeChartRef = ref()

const colors = [
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]

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
    data: [1,2,3],
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
      data: [1,2,3],
      type: 'line',
      itemStyle: {
        color: 'rgb(46,211,111)'
      }
    }
  ]
})

function loadResponseTimeData() {
  const timeStamp = [1,2,3]
  const responseData = [1,2,3]
  responseTimeChartRef.value.setOption({ 
      txAxis: [
      {
        data: timeStamp
      }
    ],
    series: [
      {
        data: responseData
      },
    ]
    })
}

defineExpose({ loadResponseTimeData })

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
      <el-row>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="percentRecord" :color="colors">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">录制</span>
            </template>
          </el-progress>
        </el-col>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="percentDataHandle" :color="colors">
            <template #default="{ percentage }">
              <span class="percentage-value">{{ percentage }}%</span>
              <span class="percentage-label">数据处理</span>
            </template>
          </el-progress>
        </el-col>
        <el-col :span="8">
          <el-progress type="dashboard" :percentage="percentDataAnalyse" :color="colors">
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
    <div style="width:300px;height:200px">
      <v-chart class="chart" ref="responseTimeChartRef" :option="responseTimeData"  />
    </div>
    

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

</style>