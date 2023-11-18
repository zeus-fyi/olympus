import os
from datetime import timezone

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

# self.id = id
# self.peer_id = peer_id
# self.date = date
# self.message = message
# self.out = out
# self.mentioned = mentioned
# self.media_unread = media_unread
# self.silent = silent
# self.post = post
# self.from_scheduled = from_scheduled
# self.legacy = legacy
# self.edit_hide = edit_hide
# self.pinned = pinned
# self.noforwards = noforwards
# self.invert_media = invert_media
# self.from_id = from_id
# self.fwd_from = fwd_from
# self.via_bot_id = via_bot_id
# self.reply_to = reply_to
# self.media = media
# self.reply_markup = reply_markup
# self.entities = entities
# self.views = views
# self.forwards = forwards
# self.replies = replies
# self.edit_date = edit_date
# self.post_author = post_author
# self.grouped_id = grouped_id
# self.reactions = reactions
# self.restriction_reason = restriction_reason
# self.ttl_period = ttl_period


async def start_client(token_in, group):
    client = TelegramClient('session_name', API_ID, API_HASH)
    await client.start(phone=lambda: get_number(), code_callback=lambda: token_in)
    await client.connect()
    msgs = []
    async for dialog in client.iter_dialogs():
        if dialog.is_group and str.startswith(dialog.name, group):
            print(f'Group name: {dialog.name}')
            messages = await client.get_messages(dialog, limit=50)
            for message in messages:
                if message.text is None or message.text == '':
                    continue
                # print(f'Message from {message.sender_id}: {message.text}')
                # # convert datetime object to unix timestamp
                unix_timestamp = message.date.replace(tzinfo=timezone.utc).timestamp()
                msgs.append({
                    'timestamp': unix_timestamp,
                    'group_name': dialog.name,
                    'sender_id': message.sender_id,
                    'message_text': message.text,
                    'is_reply': message.is_reply,
                    'is_channel': message.is_channel,
                    'is_group': message.is_group,
                    'is_private': message.is_private,
                    'first_name': message.sender.first_name,
                    'last_name': message.sender.last_name,
                    'phone': message.sender.phone,
                    'mutual_contact': message.sender.mutual_contact,
                    'username': message.sender.username,
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

