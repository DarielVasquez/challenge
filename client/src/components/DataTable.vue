<template>
  <div class="container mx-auto">
    Search:
    <input
      type="text"
      ref="inputRef"
      v-model.trim="searchQuery"
      @keyup.enter="changeQuery, getEmails()"
      class="border border-gray-500 rounded-md bg-white px-2 py-1"
    />
    <div v-if="errorMsg">{{ errorMsg }}</div>
    <div class="lg:grid lg:grid-cols-3 bg-white">
      <div class="flex flex-col lg:col-span-2">
        <h1 class="text-2xl font-bold mb-4">Emails Data Table</h1>
        <div class="overflow-x-auto">
          <table class="table-auto w-full">
            <thead>
              <tr class="bg-gray-200 cursor-pointer">
                <th class="px-4 py-2 w-2/6">Subject</th>
                <th class="px-4 py-2 w-1/6">Origin</th>
                <th class="px-4 py-2 w-1/6">From</th>
                <th class="px-4 py-2 w-1/6">To</th>
                <th class="px-4 py-2 w-1/6">Folder</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="email in emails"
                :key="email._timestamp"
                class="border-t mx-auto cursor-pointer"
                @click="selectFeaturedEmail(email)"
              >
                <td class="px-4 py-2 w-2/6">
                  <div v-html="highlight(email.subject)"></div>
                </td>
                <td class="px-4 py-2 w-1/6" v-html="highlight(email.x_origin)"></td>
                <td class="px-4 py-2 w-1/6" v-html="highlight(email.from)"></td>
                <td class="px-4 py-2 w-1/6" v-html="highlight(email.to_0_)"></td>
                <td class="px-4 py-2 w-1/6" v-html="getEmailFolder(email.x_folder)"></td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <div class="px-5 py-10 overflow-x-auto">
        <div v-html="decodedContent(highlight(featuredEmail.content))"></div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import he from 'he'
export default {
  name: 'DataTable',
  data() {
    return {
      emails: [],
      featuredEmail: {},
      searchQuery: '',
      searchQueryASCII: '',
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
            this.emails = response.data.hits
            this.featuredEmail = response.data.hits[0]
            this.errorMsg = ''
          }
        })
        .catch((error) => {
          this.errorMsg = 'Error retrieving data.'
        })
    },
    getEmailFolder(emailFolder) {
      // Split to get the file folder from the string
      const folder = emailFolder || ''
      return folder.split('\\').pop()
    },
    decodedContent(content) {
      if (content !== undefined) {
        // Use he.decode() to decode HTML entities
        const decoded = he.decode(content)
        // Replace newline characters with HTML line breaks
        return decoded.replace(/\n/g, '<br>')
      }
    },
    highlight(fileName) {
      if (this.searchQuery === '') {
        return fileName
      } else {
        return fileName?.replace(
          this.searchQueryASCII,
          `<span class="bg-yellow-500">${this.searchQueryASCII}</span>`
        )
      }
    },
    selectFeaturedEmail(email) {
      this.featuredEmail = email
    }
  },
  computed: {
    changeQuery() {
      // Update sql query with the search input
      this.payload.query.sql = `SELECT * FROM emails WHERE match_all('${this.searchQueryASCII}')`
    }
  }
}
</script>

<style lang="scss" scoped></style>
