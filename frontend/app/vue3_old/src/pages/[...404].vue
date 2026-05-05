<script setup lang="ts">

import {navigateTo} from "@/router";

import image404 from "@/assets/images/exception/404.svg";

const toHome = async () => {
  await navigateTo('/');
}
</script>

<template>
  <div class="notfound">
    <n-image
      :src="image404"
      alt=""
      class="fullscreen-image"
      @click="toHome"
      preview-disabled
    />
  </div>
</template>

<style scoped>
.notfound {
  @apply m-auto flex items-center justify-center text-9xl font-black w-screen h-screen;
  @apply bg-gradient-to-r from-green-400 to-blue-500;

  text-shadow: 10px 5px 15px slategrey;
}

.fullscreen-image {
  width: 100%;
  height: 100%;
  object-fit: fill; /* Ensures the image covers the entire area */
}
</style>

<route lang="yaml">
meta:
layout: 404
</route>
