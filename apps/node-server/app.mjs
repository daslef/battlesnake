import path from 'node:path'
import { spawn } from 'node:child_process'

import fastify from 'fastify'
import fastifyCors from '@fastify/cors'

const pathToExecutable = path.join(import.meta.dirname, '..', 'rules', 'battlesnake.exe')

const app = fastify()

app.register(fastifyCors, {
    origin: '*'
})

app.get('/new', (request, response) => {
    const child = spawn(pathToExecutable, ['play', '--name', 'ithub-starter-js', '--url', 'https://snapepy.onrender.com', '--name', 'ithub-starter-js', '--url', 'https://snapepy.onrender.com'])

    response.send()
    
    // child.stdout.on('data', data => {
    // console.log(JSON.parse(data.toString()))
    // events.push(JSON.parse(data.toString()).board)
    // })

    child.stderr.on('data', error => {
        console.error(error.toString())
    });

    child.on('close', (code) => {
        console.log(code)
    });
})

app.listen({ host: 'localhost', port: 5001 })