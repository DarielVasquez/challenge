<template>
  <div class="flex justify-around p-5">
    <div>
      <v-icon
        :class="[chevronStyles, { 'opacity-25 cursor-default': currentPage === 1 }]"
        :disabled="currentPage > 1"
        @click="firstPage"
        name="fa-angle-double-left"
      ></v-icon>
      <v-icon
        :class="[chevronStyles, { 'opacity-25 cursor-default': currentPage === 1 }]"
        :disabled="currentPage > 1"
        @click="prevPage"
        name="fa-chevron-left"
      ></v-icon>
    </div>
    <div>{{ Math.ceil(currentPage) }} of {{ totalPages }}</div>
    <div>
      <v-icon
        :class="[chevronStyles, { 'opacity-25 cursor-default': currentPage === totalPages }]"
        :disabled="currentPage < totalPages"
        @click="nextPage"
        name="fa-chevron-right"
      ></v-icon>
      <v-icon
        :class="[chevronStyles, { 'opacity-25 cursor-default': currentPage === totalPages }]"
        :disabled="currentPage < totalPages"
        @click="lastPage"
        name="fa-angle-double-right"
      ></v-icon>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
export default {
  props: ['currentPage', 'totalPages'],
  setup(props, { emit }) {
    const firstPage = () => {
      if (props.currentPage > 1) {
        emit('update:currentPage', 1)
      }
    }

    const prevPage = () => {
      if (props.currentPage > 1) {
        emit('update:currentPage', props.currentPage - 1)
      }
    }

    const nextPage = () => {
      if (props.currentPage < props.totalPages) {
        emit('update:currentPage', props.currentPage + 1)
      }
    }

    const lastPage = () => {
      if (props.currentPage < props.totalPages) {
        emit('update:currentPage', props.totalPages)
      }
    }

    const chevronStyles = computed(() => {
      return 'text-slate-600 cursor-pointer'
    })

    return {
      firstPage,
      prevPage,
      nextPage,
      lastPage,
      chevronStyles
    }
  }
}
</script>

<style lang="scss" scoped></style>
