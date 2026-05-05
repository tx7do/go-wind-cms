<template>
  <n-dialog-provider>
    <n-notification-provider>
      <n-message-provider>
        <n-loading-bar-provider>
          <n-modal-provider>
            <slot name="default"></slot>
          </n-modal-provider>
        </n-loading-bar-provider>
      </n-message-provider>
    </n-notification-provider>
  </n-dialog-provider>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {
  NDialogProvider,
  NNotificationProvider,
  NMessageProvider,
  NLoadingBarProvider
} from 'naive-ui'

export default defineComponent({
  name: 'ApplicationProvider',
  components: {
    NDialogProvider,
    NNotificationProvider,
    NMessageProvider,
    NLoadingBarProvider
  },
  setup() {
    return {}
  }
})
</script>
