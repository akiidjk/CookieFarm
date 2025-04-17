<script setup lang="ts">
import { ref } from 'vue'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { HOST } from '@/lib/config'
import { Icon } from '@iconify/vue'

const password = ref('')
const error = ref('')
const typeInput = ref<'password' | 'text'>('password')
const token = useCookie('token', {
    maxAge: 60 * 60 * 24 * 7,
    path: '/',
    secure: true,
    // httpOnly: true
})

async function handleFormSubmit() {
    try {
        const res = await $fetch<{ token: string }>(HOST + '/api/v1/auth/login', {
            method: 'POST',
            body: { password: password.value },
        })
        token.value = res.token
        error.value = ''
        console.log('Token:', token.value)
        await navigateTo('/')
    } catch (err: unknown) {
        console.error(err)
        const errorObj = err as { data?: { message?: string } }
        error.value = errorObj?.data?.message || 'Login failed'
    }
}

function togglePasswordVisibility() {
    typeInput.value = typeInput.value === 'password' ? 'text' : 'password'
}
</script>

<template>
    <div class="w-full h-screen flex items-center justify-center px-4">
        <Card class="w-full max-w-sm">
            <CardHeader>
                <div class="flex flex-col items-center justify-center">
                    <img src="/assets/logo/logo.png" alt="Logo" class="w-16 h-16 mb-2">
                    <h1 class="text-2xl font-bold">Cookie Farm</h1>
                </div>
                <CardTitle class="text-2xl">Login</CardTitle>
                <CardDescription>
                    Enter the password below to login.
                </CardDescription>
            </CardHeader>
            <CardContent class="grid gap-4">
                <div class="grid gap-2">
                    <Label for="password">Password</Label>
                    <div class="flex gap-2">
                        <Input id="password" v-model="password" :type="typeInput" required
                            placeholder="Enter password" />
                        <Button class="w-10 h-10" type="button"
                            :aria-label="typeInput === 'password' ? 'Show password' : 'Hide password'"
                            @click="togglePasswordVisibility">
                            <Icon
                                :icon="typeInput === 'password' ? 'material-symbols:visibility-off' : 'material-symbols:visibility'"
                                class="w-5 h-5" />
                        </Button>
                    </div>
                </div>
                <p v-if="error" class="text-sm text-red-500">{{ error }}</p>
            </CardContent>
            <CardFooter>
                <Button class="w-full" @click="handleFormSubmit">
                    Sign in
                </Button>
            </CardFooter>
        </Card>
    </div>
</template>
