<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useNavigationStore } from '@/stores'

const navigationStore = useNavigationStore()
const result = ref<any>(null)
const error = ref<any>(null)

async function testNavigation() {
  try {
    console.log('Testing navigation API...')
    const res = await navigationStore.listNavigation(
      { page: 1, pageSize: 10 },
      { location: 'HEADER', isActive: true }
    )
    console.log('Response:', res)
    result.value = res
  } catch (e) {
    console.error('Error:', e)
    error.value = e
  }
}

onMounted(() => {
  testNavigation()
})
</script>

<template>
  <div style="padding: 20px">
    <h1>Navigation API Test</h1>
    <n-button @click="testNavigation">Test Navigation</n-button>

    <h2>Result:</h2>
    <pre>{{ JSON.stringify(result, null, 2) }}</pre>

    <h2>Error:</h2>
    <pre>{{ JSON.stringify(error, null, 2) }}</pre>
  </div>
</template>

