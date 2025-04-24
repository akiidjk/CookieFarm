<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import {
    useVueTable,
    createColumnHelper,
    getCoreRowModel,
} from '@tanstack/vue-table'
import FormConfig from '~/components/form-config/FormConfig.vue'

interface Flag {
    submit_time: number
    response_time: number
    service_port: number
    team_id: number
    id: string
    flag_code: string
    service_name: string
    status: string
}

interface FlagsResponse {
    n_flags: number
    flags: Flag[]
}

const token = useCookie('token')

const page = ref(0)
const pageSize = 10

const flagsResponse = ref<FlagsResponse | null>(null)
const isLoading = ref(true)

const fetchFlags = async () => {
    isLoading.value = true
    try {
        const res = await $fetch<FlagsResponse>(
            `http://localhost:8080/api/v1/flags/${pageSize}?offset=${page.value * pageSize}`,
            {
                headers: {
                    Authorization: 'Bearer ' + token.value!
                }
            }
        )
        flagsResponse.value = res
    } catch (err) {
        console.error('Errore caricamento flags:', err)
    } finally {
        isLoading.value = false
    }
}

watch(page, fetchFlags, { immediate: true })

const flags = computed(() => flagsResponse.value?.flags ?? [])
const totalFlags = computed(() => flagsResponse.value?.n_flags ?? 0)
const totalPages = computed(() => Math.ceil(totalFlags.value / pageSize))

const columnHelper = createColumnHelper<Flag>()

const columns = [
    columnHelper.accessor('flag_code', { header: 'Flag' }),
    columnHelper.accessor('status', { header: 'Status' }),
    columnHelper.accessor('service_name', { header: 'Service' }),
    columnHelper.accessor('submit_time', {
        header: 'Submit Time',
        cell: info => new Date(info.getValue() * 1000).toLocaleString()
    }),
    columnHelper.accessor('response_time', {
        header: 'Response Time',
        cell: info => new Date(info.getValue() * 1000).toLocaleString()
    }),
    columnHelper.accessor('team_id', { header: 'Team' }),
    columnHelper.accessor('service_port', { header: 'Port' }),
]

const table = useVueTable({
    data: flags,
    columns,
    getCoreRowModel: getCoreRowModel()
})
</script>


<template>
    <div>
        <FormConfig />
        <h1 class="text-2xl font-bold mb-2">Welcome to Cookie Farm</h1>
        <p class="mb-4 text-gray-600">Gestione delle flag con caricamento paginato lato server</p>

    </div>
</template>
