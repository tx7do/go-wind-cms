import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Tenant {
  id: string
  name: string
  description?: string
  logo?: string
  owner_id: string
  created_at: string
  updated_at: string
  status: 'active' | 'inactive'
}

export const useTenantStore = defineStore('tenant', () => {
  // State
  const currentTenant = ref<Tenant | null>(null)
  const tenants = ref<Tenant[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)

  // Computed
  const currentTenantId = computed(() => currentTenant.value?.id)
  const isLoading = computed(() => loading.value)
  const hasError = computed(() => error.value !== null)

  // Actions
  async function fetchTenants() {
    loading.value = true
    error.value = null
    try {
      // TODO: Replace with actual API call
      // const response = await api.getTenants()
      // tenants.value = response.data

      // Mock data for now
      tenants.value = [
        {
          id: 'tenant-1',
          name: 'Acme Corp',
          description: 'Main company tenant',
          logo: '/logo.png',
          owner_id: 'user-1',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          status: 'active',
        },
        {
          id: 'tenant-2',
          name: 'Tech Startup',
          description: 'Technology startup',
          logo: '/logo.png',
          owner_id: 'user-1',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          status: 'active',
        },
        {
          id: 'tenant-3',
          name: 'Creative Agency',
          description: 'Creative design agency',
          logo: '/logo.png',
          owner_id: 'user-1',
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString(),
          status: 'active',
        },
      ]

      // Set first tenant as current if none selected
      if (!currentTenant.value && tenants.value.length > 0) {
        currentTenant.value = tenants.value[0]
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to fetch tenants:', e)
    } finally {
      loading.value = false
    }
  }

  async function switchTenant(tenantId: string) {
    const tenant = tenants.value.find(t => t.id === tenantId)
    if (tenant) {
      currentTenant.value = tenant
      // TODO: Update API headers with new tenant ID
      // TODO: Reload data for new tenant
    } else {
      console.error(`Tenant ${tenantId} not found`)
    }
  }

  async function createTenant(data: Omit<Tenant, 'id' | 'created_at' | 'updated_at'>) {
    loading.value = true
    error.value = null
    try {
      // TODO: Replace with actual API call
      // const response = await api.createTenant(data)
      // const newTenant = response.data

      // Mock for now
      const newTenant: Tenant = {
        ...data,
        id: `tenant-${Date.now()}`,
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
      }

      tenants.value.push(newTenant)
      return newTenant
    } catch (e) {
      error.value = e as Error
      console.error('Failed to create tenant:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function updateTenant(tenantId: string, data: Partial<Tenant>) {
    loading.value = true
    error.value = null
    try {
      // TODO: Replace with actual API call
      // const response = await api.updateTenant(tenantId, data)

      const index = tenants.value.findIndex(t => t.id === tenantId)
      if (index !== -1) {
        tenants.value[index] = {
          ...tenants.value[index],
          ...data,
          updated_at: new Date().toISOString(),
        }

        // Update current tenant if it's the one being updated
        if (currentTenant.value?.id === tenantId) {
          currentTenant.value = tenants.value[index]
        }
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to update tenant:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  async function deleteTenant(tenantId: string) {
    loading.value = true
    error.value = null
    try {
      // TODO: Replace with actual API call
      // await api.deleteTenant(tenantId)

      const index = tenants.value.findIndex(t => t.id === tenantId)
      if (index !== -1) {
        tenants.value.splice(index, 1)

        // Switch to another tenant if current one was deleted
        if (currentTenant.value?.id === tenantId && tenants.value.length > 0) {
          currentTenant.value = tenants.value[0]
        } else if (tenants.value.length === 0) {
          currentTenant.value = null
        }
      }
    } catch (e) {
      error.value = e as Error
      console.error('Failed to delete tenant:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  function $reset() {
    currentTenant.value = null
    tenants.value = []
    loading.value = false
    error.value = null
  }

  return {
    // State
    currentTenant,
    tenants,
    loading,
    error,

    // Computed
    currentTenantId,
    isLoading,
    hasError,

    // Actions
    fetchTenants,
    switchTenant,
    createTenant,
    updateTenant,
    deleteTenant,
    $reset,
  }
}, {
  persist: {
    key: 'tenant-store',
    storage: localStorage,
    paths: ['currentTenant'],
  },
})

