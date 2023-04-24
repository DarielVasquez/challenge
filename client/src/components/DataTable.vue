<template>
  <div class="container mx-auto">
    <div v-if="errorMsg">{{ errorMsg }}</div>
    <div class="lg:grid lg:grid-cols-3 bg-white">
      <div class="flex flex-col lg:col-span-2">
        <h1 class="text-2xl font-bold mb-4">Emails Data Table</h1>
        <div class="flex justify-between p-2">
          <div>
            <span>Show</span>
            <select class="text-center" name="limit" id="limit" @change="handleLimitChange">
              <option :value="limit" v-for="(limit, index) in limits" :key="index">
                {{ limit }}
              </option>
            </select>
            <span>entries</span>
          </div>
          <div>
            Search:
            <input
              type="text"
              ref="inputRef"
              v-model.trim="searchQuery"
              @keyup.enter="changeQuery, getEmails()"
              class="border border-gray-500 rounded-md bg-white px-2 py-1"
            />
          </div>
        </div>
        <div class="overflow-x-auto">
          <table class="table-auto w-full">
            <thead>
              <tr class="bg-gray-200 cursor-pointer">
                <th class="px-4 py-2 min-w-[120px]" @click="sortTable('subject')">
                  Subject
                  <sort-icons-vue :sortDirection="sortDirection" header="subject"></sort-icons-vue>
                </th>
                <th class="px-4 py-2 min-w-[105px]" @click="sortTable('x_origin')">
                  Origin
                  <sort-icons-vue :sortDirection="sortDirection" header="x_origin"></sort-icons-vue>
                </th>
                <th class="px-4 py-2 min-w-[100px]" @click="sortTable('from')">
                  From
                  <sort-icons-vue :sortDirection="sortDirection" header="from"></sort-icons-vue>
                </th>
                <th class="px-4 py-2 min-w-[100px]" @click="sortTable('to_0_')">
                  To
                  <sort-icons-vue :sortDirection="sortDirection" header="to_0_"></sort-icons-vue>
                </th>
                <th class="px-4 py-2 min-w-[105px]" @click="sortTable('x_folder')">
                  Folder
                  <sort-icons-vue :sortDirection="sortDirection" header="x_folder"></sort-icons-vue>
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="email in sortedTableEmails"
                :key="email._timestamp"
                class="border-t mx-auto cursor-pointer"
                @click="selectFeaturedEmail(email)"
              >
                <td class="px-4 py-2" v-html="highlight(email.subject)"></td>
                <td class="px-4 py-2" v-html="highlight(email.x_origin)"></td>
                <td class="px-4 py-2" v-html="highlight(email.from)"></td>
                <td class="px-4 py-2" v-html="highlight(email.to_0_)?.split(',').splice(0, 5)"></td>
                <td class="px-4 py-2" v-html="highlight(email.x_folder)"></td>
              </tr>
            </tbody>
          </table>
        </div>
        <pagination-vue
          :currentPage="currentPage"
          :totalPages="totalPages"
          @update:currentPage="updateCurrentPage"
        ></pagination-vue>
        <div class="flex justify-center">
          Showing {{ indexOfFirstEmail + 1 }} to
          {{ indexOfLastEmail > emails.length ? emails.length : indexOfLastEmail }} of
          {{ emails.length }} entries in {{ elapsedTime.toFixed(2) }}ms
        </div>
      </div>
      <div class="px-5 py-10 overflow-x-auto">
        <div v-html="highlight(featuredEmail.content)"></div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import he from 'he'
import SortIconsVue from './SortIcons.vue'
import PaginationVue from './Pagination.vue'
export default {
  name: 'DataTable',
  components: {
    SortIconsVue,
    PaginationVue
  },
  data() {
    return {
      emails: [],
      featuredEmail: {}, // Email that is displayed on-screen
      sortBy: '', // Current sort column
      sortDirection: {}, // Sort direction for table headers
      currentPage: 1, // Current page for pagination in data table
      emailsPerPage: 5,
      limits: [5, 10, 15, 20, 25, 50, 100], // Array for amount of entries per page
      elapsedTime: 0, // Elapsed time it takes for the api to execute
      searchQuery: '',
      searchQueryASCII: '', // Search query converted to ASCII for api reading
      errorMsg: '',
      payload: {
        query: {
          sql: 'SELECT * FROM emails',
          size: 100
        }
      }
    }
  },
  created() {
    this.getEmails()
  },
  mounted() {
    this.$refs.inputRef.focus()
  },
  watch: {
    // Update search query to replace single quotes to ASCII
    searchQuery(search) {
      this.searchQueryASCII = search.replace(/'/g, '&#39;')
    }
  },
  methods: {
    getEmails() {
      const startTime = performance.now()
      axios
        .post('http://localhost:5080/api/default/_search', this.payload, {
          auth: {
            username: 'root@example.com',
            password: 'cQhEVBboURHqqKd6'
          }
        })
        .then((response) => {
          if (response.data.hits.length === 0) {
            this.errorMsg = 'No data was found, please try a different query.'
          } else {
            this.emails = response.data.hits.map((email) => {
              const endTime = performance.now()
              this.elapsedTime = endTime - startTime
              email.content = he.decode(email.content).replace(/\n/g, '<br>')
              email.subject = he.decode(email.subject)
              email.x_folder = email.x_folder.split('\\').pop()
              this.currentPage = 1 // Set current page back to 1
              for (const key in this.sortDirection) {
                // Set all sort headers to default
                if (this.sortDirection.hasOwnProperty(key)) {
                  this.sortDirection[key] = ''
                }
              }
              return email
            })
            this.featuredEmail = response.data.hits[0]
            this.errorMsg = ''
          }
        })
        .catch((error) => {
          this.errorMsg = 'Error retrieving data.'
        })
    },
    highlight(fileName) {
      if (this.searchQuery === '') {
        return fileName
      } else {
        const searchPattern = new RegExp(`(?<!<[^>]*)(${this.searchQuery})`, 'gi') // Excludes <br> tags from getting highlighted
        return fileName?.replaceAll(searchPattern, (match) => {
          return `<span class="bg-yellow-500">${match}</span>`
        })
      }
    },
    selectFeaturedEmail(email) {
      this.featuredEmail = email
    },
    sortTable(column) {
      if (this.sortBy === column) {
        this.sortDirection[column] = this.sortDirection[column] === 'asc' ? 'desc' : 'asc'
      } else {
        this.sortBy = column
        this.sortDirection[column] = 'asc'
      }
      for (const key in this.sortDirection) {
        if (this.sortDirection.hasOwnProperty(key)) {
          if (column !== key) {
            this.sortDirection[key] = ''
          }
        }
      }
    },
    handleLimitChange(event) {
      this.currentPage = 1
      this.emailsPerPage = event.target.value
    },
    updateCurrentPage(value) {
      this.currentPage = value
    }
  },
  computed: {
    changeQuery() {
      // Update sql query with the search input
      this.payload.query.sql = `SELECT * FROM emails WHERE match_all('${this.searchQueryASCII}')`
    },
    sortedTableEmails() {
      // Sort the array of emails according to the table header
      const column = this.sortBy
      const direction = this.sortDirection[column]
      return this.filteredEmails.sort((a, b) => {
        const aValue = a[column]?.toUpperCase()
        const bValue = b[column]?.toUpperCase()
        if (aValue < bValue) return direction === 'asc' ? -1 : 1
        if (aValue > bValue) return direction === 'asc' ? 1 : -1
        return 0
      })
    },
    indexOfLastEmail() {
      return this.currentPage * Math.round(this.emailsPerPage / 5) * 5
    },
    indexOfFirstEmail() {
      return this.indexOfLastEmail - Math.round(this.emailsPerPage / 5) * 5
    },
    filteredEmails() {
      return this.emails.slice(this.indexOfFirstEmail, this.indexOfLastEmail)
    },
    totalPages() {
      return Math.ceil(this.emails.length / this.emailsPerPage)
    }
  }
}
</script>

<style lang="scss" scoped></style>
