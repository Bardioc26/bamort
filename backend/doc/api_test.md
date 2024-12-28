/*
API Testing

    Register User:
        POST /register
        Body: { "username": "test", "password": "12345", "email": "test@example.com" }

    Login User:
        POST /login
        Body: { "username": "test", "password": "12345" }

    Create Character:
        POST /characters
        Body: { "user_id": 1, "name": "Hero", "rasse": "Elf", "typ": "Warrior", "age": 25 }

    Get Characters:
        GET /characters
*/