<script setup lang="ts">
import NaiveProvider from "~/components/NaiveProvider.vue";
</script>

<template>
  <n-loading-bar-provider>
    <n-dialog-provider>
      <n-notification-provider>
        <n-message-provider>
          <NaiveProvider>
            <NuxtLayout>
              <NuxtPage/>
            </NuxtLayout>
          </NaiveProvider>
        </n-message-provider>
      </n-notification-provider>
    </n-dialog-provider>
  </n-loading-bar-provider>
</template>
