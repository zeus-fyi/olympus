import os

from flask import Flask, jsonify, request
from telethon import TelegramClient

app = Flask(__name__)


def get_number():
    return "17575828406"


@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({'status': 'running'})


SESSION = 'session_name'
API_ID = int(os.environ['TELEGRAM_API_ID'])
API_HASH = os.environ['TELEGRAM_API_HASH']


async def start_client(token_in, group):
    client = TelegramClient('session_name', API_ID, API_HASH)
    await client.start(phone=lambda: get_number(), code_callback=lambda: token_in)
    await client.connect()
    msgs = []
    async for dialog in client.iter_dialogs():
        if dialog.is_group and str.startswith(dialog.name, group):
            print(f'Group name: {dialog.name}')

            messages = await client.get_messages(dialog, limit=10)
            for message in messages:
                print(f'Message from {message.sender_id}: {message.text}')
                msgs.append({
                    'sender_id': message.sender_id,
                    'message_text': message.text,
                })
    client.disconnect()
    return msgs


@app.route('/msgs', methods=['POST'])
async def initialize_telegram_client_endpoint():
    token = request.json.get('token')
    if not token:
        return jsonify({'error': 'Token is required'}), 400
    group = request.json.get('group')
    if not group:
        return jsonify({'error': 'Group is required'}), 400
    return await start_client(token, group)


if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=8000)

