<script lang="ts" setup>
// import HelloWorld from './components/HelloWorld.vue'
import AsideView from './components/AsideView.vue';
import MainContent from './components/MainContent.vue';
</script>

<template>
  <div class="common-layout">
    <el-container>
      <!-- <el-header>Header</el-header> -->
      <el-container>
        <el-aside width="200px">
          <AsideView/>
        </el-aside>
        <el-main>
          <MainContent></MainContent>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<style>

</style>
