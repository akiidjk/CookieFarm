<template>
    <NavigationMenu>
        <NavigationMenuList>
            <NavigationMenuItem>
                <NavigationMenuLink>
                    <NuxtLink to="/">Home</NuxtLink>
                </NavigationMenuLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
                <NavigationMenuLink>
                    <NuxtLink to="/events">Events</NuxtLink>
                </NavigationMenuLink>
            </NavigationMenuItem>
        </NavigationMenuList>
    </NavigationMenu>
    <NuxtLayout>
        <NuxtPage />
    </NuxtLayout>
</template>
