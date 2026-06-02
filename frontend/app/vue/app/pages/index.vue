<script setup lang="ts">
// Home page with scroll reveal
onMounted(() => {
  const intersectionObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          entry.target.classList.add('is-visible')
          intersectionObserver.unobserve(entry.target)
        }
      })
    },
    { threshold: 0.1, rootMargin: '0px 0px -100px 0px' },
  )

  const observeElements = (root: Element | Document = document) => {
    root.querySelectorAll('.scroll-reveal-item:not(.is-visible)').forEach((el) => {
      intersectionObserver.observe(el)
    })
  }
  observeElements()

  const mutationObserver = new MutationObserver((mutations) => {
    for (const mutation of mutations) {
      mutation.addedNodes.forEach((node) => {
        if (node.nodeType !== Node.ELEMENT_NODE) return
        const el = node as Element
        if (el.classList?.contains('scroll-reveal-item') && !el.classList.contains('is-visible')) {
          intersectionObserver.observe(el)
        }
        observeElements(el)
      })
    }
  })
  mutationObserver.observe(document.body, { childList: true, subtree: true })

  onUnmounted(() => {
    intersectionObserver.disconnect()
    mutationObserver.disconnect()
  })
})
</script>

<template>
  <div class="w-full">
    <HomeHeroSection />
    <HomeFeaturedPostsSection />
    <HomeCategoryListSection />
    <HomePopularTagsSection />
    <HomeLatestPostsSection />
    <HomeFeaturesSection />
  </div>
</template>
