<script setup lang="ts">
const props = withDefaults(defineProps<{
  src?: string
  alt?: string
  fallback?: string
}>(), {
  alt: '',
  fallback: '/placeholder.png',
})

const imgSrc = ref(props.src || props.fallback)

function handleError(event: Event) {
  const img = event.target as HTMLImageElement
  if (img.src.endsWith(props.fallback)) return
  imgSrc.value = props.fallback
}

watch(() => props.src, (newSrc) => {
  imgSrc.value = newSrc || props.fallback
})
</script>

<template>
  <img :src="imgSrc" :alt="alt" @error="handleError" />
</template>
