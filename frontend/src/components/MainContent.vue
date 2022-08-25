<script setup lang="ts">
import {reactive, ref, inject, Ref, onMounted, computed} from 'vue'

const processingStatus = ref(2)

const countNumber = ref(3)
const percentRecord = ref(60)
const percentDataHandle = ref(60)
const percentDataAnalyse = ref(20)


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