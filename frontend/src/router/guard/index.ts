import type { Router } from 'vue-router'
import NProgress from 'nprogress' // progress bar
import { setRouteEmitter } from '@/utils/route-listener'

// const isLoading = ref(false)
function setupPageGuard(router: Router) {
  router.beforeEach(async (to) => {
    // emit route change
    setRouteEmitter(to)
    NProgress.start()
    // isLoading.value = true
  })

  router.afterEach(async () => {
    NProgress.done()
    // isLoading.value = false
  })
}
// export { isLoading }
export default function createRouteGuard(router: Router) {
  setupPageGuard(router)
}
