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
export default {
  props: ['currentPage', 'totalPages'],
  methods: {
    firstPage() {
      if (this.currentPage > 1) {
        this.$emit('update:currentPage', 1)
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.$emit('update:currentPage', this.currentPage - 1)
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.$emit('update:currentPage', this.currentPage + 1)
      }
    },
    lastPage() {
      if (this.currentPage < this.totalPages) {
        this.$emit('update:currentPage', this.totalPages)
      }
    }
  },
  computed: {
    // Styles
    chevronStyles() {
      return 'text-slate-600 cursor-pointer'
    }
  }
}
</script>

<style lang="scss" scoped></style>
