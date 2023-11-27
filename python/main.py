from flask import Flask, request, jsonify
from database import cursor
from datetime import datetime
from todo import Todo


app = Flask(__name__)


@app.route('/todos', methods=['GET'])
def get_todos():
    cursor.execute("SELECT * FROM todos")
    todos = cursor.fetchall()
    return jsonify(todos)


@app.route('/todos', methods=['POST'])
def create_todo():
    data = request.get_json()
    title = data['title']
    description = data['description']
    added_date = datetime.now()

    cursor.execute("INSERT INTO todos (title, description, added_date) VALUES (%s, %s, %s) RETURNING id",
                   (title, description, added_date))
    todo_id = cursor.fetchone()['id']

    todo = Todo(id=todo_id, title=title, description=description,
                added_date=added_date, completed_date=None)

    return jsonify(todo.__dict__)


@app.route('/todos/<int:todo_id>', methods=['GET'])
def get_todo(todo_id):
    cursor.execute("SELECT * FROM todos WHERE id = %s", (todo_id,))
    todo = cursor.fetchone()
    return jsonify(todo)

# Route to update a todo


@app.route('/todos/<int:todo_id>', methods=['PUT'])
def update_todo(todo_id):
    data = request.get_json()
    title = data['title']
    description = data['description']
    completed_date = data.get('completed_date', None)

    cursor.execute("UPDATE todos SET title = %s, description = %s, completed_date = %s WHERE id = %s",
                   (title, description, completed_date, todo_id))

    cursor.execute("SELECT * FROM todos WHERE id = %s", (todo_id,))
    updated_todo = cursor.fetchone()

    return jsonify(updated_todo)


@app.route('/todos/<int:todo_id>', methods=['DELETE'])
def delete_todo(todo_id):
    cursor.execute("DELETE FROM todos WHERE id = %s", (todo_id,))
    return '', 204


if __name__ == '__main__':
    app.run(port=8080)
