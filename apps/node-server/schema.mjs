const schema = {
    type: 'object',
    properties: {
        snakes: {
            type: 'array',
            items: {
                type: 'object',
                required: ['snakeName', 'snakeUrl'],
                properties: {
                    snakeName: {
                        type: 'string'
                    },
                    snakeUrl: {
                        type: 'string'
                    }
                }
            }
        },
        field: {
            type: 'object',
            required: ['width', 'height'],
            properties: {
                width: {
                    type: 'number'
                },
                height: {
                    type: 'number'
                }
            }
        }
    }
}

export default schema