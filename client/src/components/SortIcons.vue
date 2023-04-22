<template>
  <span>
    <v-icon v-if="!sortDirection[header] || sortDirection[header] === ''" name="fa-sort" />
    <v-icon v-if="sortDirection[header] === 'asc'" name="fa-sort-up" />
    <v-icon v-if="sortDirection[header] === 'desc'" name="fa-sort-down" />
  </span>
</template>

<script>
export default {
  props: ['sortDirection', 'header']
}
</script>

<style lang="scss" scoped></style>
