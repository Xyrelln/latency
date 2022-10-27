<script setup lang="ts">
import {reactive, ref, onMounted} from 'vue'

import { GetCurrentVersion } from '../../wailsjs/go/app/Api' 

const version = ref('0.2.x-dev')

const getCurrentVersion = () => {
  GetCurrentVersion().then((res: string) => {
    if ( res !== '') {
      version.value = res
    }
  }).catch(err => { console.log(err) })
}

onMounted(() => {
  getCurrentVersion()
})
</script>

<template>
  <el-scrollbar style="height: calc(100vh - 100px);width: calc(100vw - 60px)">
    <div class="describe">
      <el-row>
        <el-col :span="4">
          版权所有
        </el-col>
        <el-col :span="8">
          云天畅享 http://www.ivcloud.net/
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="4">
          版本号
        </el-col>
        <el-col :span="4">
          {{ version }}
        </el-col>
      </el-row>
    </div>
  </el-scrollbar>
  
</template>

<style>


</style>
