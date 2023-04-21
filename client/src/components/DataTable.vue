<template>
  <div>
    <input
      type="text"
      ref="inputRef"
      v-model.trim.lazy="searchQuery"
      @keyup.enter="changeQuery, getEmails()"
      class="border border-gray-500 rounded-md text-blue-500 bg-white px-2 py-1"
    />
    <div v-if="errorMsg">{{ errorMsg }}</div>
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Name</th>
          <th>Email</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="email in emails" :key="email._timestamp">
          <td>{{ email._timestamp }}</td>
          <td>{{ email.x_origin }}</td>
          <td v-html="decodedContent(email.content)"></td>
          <hr />
        </tr>
      </tbody>
    </table>
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
    decodedContent(content) {
      // Use he.decode() to decode HTML entities
      const decoded = he.decode(content)
      // Replace newline characters with HTML line breaks
      return decoded.replace(/\n/g, '<br>')
    }
  },
  computed: {
    changeQuery() {
      // Update sql query with the search input
      this.payload.query.sql = `SELECT * FROM test WHERE match_all('${this.searchQuery.replace(
        /'/g,
        '&#39;'
      )}')`
    }
  }
}
</script>

<style lang="scss" scoped></style>
