import path from 'node:path'
import { spawn } from 'node:child_process'

import fastify from 'fastify'
import fastifyCors from '@fastify/cors'

import schema from './schema.mjs'

const pathToExecutable = path.join(import.meta.dirname, '..', 'rules', 'battlesnake.exe')

const app = fastify()

app.register(fastifyCors, {
    origin: '*'
})

app.post('/new', { body: schema }, async (request, response) => {
    const childArgs = ['play']

    const { snakes, field } = request.body

    for (const snake of snakes) {
        childArgs.push('--name', snake.snakeName, '--url', snake.snakeUrl)
    }

    childArgs.push('--width', field.width, '--height', field.height)

    const child = spawn(pathToExecutable, childArgs)

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