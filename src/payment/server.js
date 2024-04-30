import app from './app.js'

const PORT = process.env.PORT || 3000 // Default to port 3000 if PORT environment variable is not set
app.listen(PORT, () => {
  console.log(`Listening on port ${PORT}`)
})
