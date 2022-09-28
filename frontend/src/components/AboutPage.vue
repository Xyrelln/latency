<script setup lang="ts">
import {reactive, ref } from 'vue'


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
          0.1.0
        </el-col>
      </el-row>
    </div>
  </el-scrollbar>
  
</template>

<style>


</style>
