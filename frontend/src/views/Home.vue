<script setup lang="ts">
import {reactive, ref, h, inject, Ref, provide, onMounted, computed, watch, onUnmounted} from 'vue'
import { ElMessage } from 'element-plus'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import Automation from '../components/Automation.vue';
import AboutPage from '../components/AboutPage.vue';

import LatencyPage from './LatencyPage.vue';

import { isWailsRun } from '@/utils/utils'



const topTabName = ref('latency')


</script>

<template>
  <el-container>
    <el-tabs type="border-card" v-model="topTabName" >
      <el-tab-pane label="延迟测试" name="latency">
        <LatencyPage></LatencyPage>
      </el-tab-pane>
      <el-tab-pane label="场景配置" name="automation">
        <Automation/>
      </el-tab-pane>
      <el-tab-pane label="关于" name="about">
        <AboutPage></AboutPage>
      </el-tab-pane>
    </el-tabs>
   
  </el-container>
    
</template>

<style scoped>

</style>
