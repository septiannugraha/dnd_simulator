// MongoDB initialization script
db = db.getSiblingDB('dnd_simulator');

// Create collections with indexes
db.createCollection('users');
db.createCollection('characters');
db.createCollection('campaigns');
db.createCollection('game_sessions');

// Create indexes for better performance
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });

db.characters.createIndex({ "user_id": 1 });
db.characters.createIndex({ "campaign_id": 1 });

db.campaigns.createIndex({ "dm_id": 1 });
db.campaigns.createIndex({ "player_ids": 1 });

db.game_sessions.createIndex({ "campaign_id": 1 });
db.game_sessions.createIndex({ "player_ids": 1 });

print('Database initialized with collections and indexes');