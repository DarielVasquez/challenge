<template>
  <div class="container mx-auto">
    <input
      type="text"
      ref="inputRef"
      v-model.trim="searchQuery"
      @keyup.enter="changeQuery, getEmails()"
      class="border border-gray-500 rounded-md text-blue-500 bg-white px-2 py-1"
    />
    <div v-if="errorMsg">{{ errorMsg }}</div>
    <div class="flex flex-col">
      <h1 class="text-2xl font-bold mb-4">Emails Data Table</h1>
      <div class="overflow-x-auto">
        <table class="table-auto w-full">
          <thead>
            <tr class="bg-gray-200">
              <th class="px-4 py-2">Subject</th>
              <th class="px-4 py-2">Origin</th>
              <th class="px-4 py-2">From</th>
              <th class="px-4 py-2">To</th>
              <th class="px-4 py-2">Folder</th>
              <th class="px-4 py-2">Content</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="email in emails" :key="email._timestamp" class="border-t mx-auto">
              <td class="px-4 py-2" v-html="highlight(email.subject)"></td>
              <td class="px-4 py-2" v-html="email.x_origin"></td>
              <td class="px-4 py-2" v-html="email.from"></td>
              <td class="px-4 py-2" v-html="email.to_0_"></td>
              <td class="px-4 py-2" v-html="getEmailFolder(email.x_folder)"></td>
              <td class="px-4 py-2">
                <div v-html="decodedContent(highlight(email.content.slice(0, 50)))"></div>
              </td>
            </tr>
          </tbody>
        </table>
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
          console.log(response.data.hits)
          this.emails = response.data.hits
          this.errorMsg = ''
        })
        .catch((error) => {
          this.errorMsg = 'Error retrieving data'
        })
    },
    getEmailFolder(emailFolder) {
      // Split to get the file folder from the string
      const folder = emailFolder || ''
      return folder.split('\\').pop()
    },
    decodedContent(content) {
      // Use he.decode() to decode HTML entities
      const decoded = he.decode(content)
      // Replace newline characters with HTML line breaks
      return decoded.replace(/\n/g, '<br>')
    },
    highlight(fileName) {
      if (this.searchQuery === '') {
        return fileName
      } else {
        return fileName.replace(
          this.searchQueryASCII,
          `<span class="bg-yellow-500">${this.searchQueryASCII}</span>`
        )
      }
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
