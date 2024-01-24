<template>
  <div class="container max-w-screen-lg">
    <div class="m-auto w-full justify-start py-4">
      <div class="flex gap-4">
        <div class="font-bold cursor-pointer" @click="() => $router.push({ name: 'main'})">Home</div>
        <div class="font-bold cursor-pointer" @click="() => $router.push({ name: 'orders'})">Orders</div>
      </div>
    </div>
  </div>
  <div class="w-full border-b-slate-300 border-b"></div>
</template>
