CREATE TABLE project_types (
    id bigserial NOT NULL primary key,
    alias varchar NOT NULL,
    name varchar NOT NULL,
    options TEXT [],
    goal_by_amount boolean NOT NULL DEFAULT FALSE,
    goal_by_people boolean NOT NULL DEFAULT FALSE,
    end_by_goal_gain boolean NOT NULL DEFAULT FALSE
);

INSERT INTO project_types (alias, name, goal_by_amount, goal_by_people, end_by_goal_gain, options) VALUES (
    'money_fast', 'Сampaign', true, false, true, ARRAY [
        'Partakers deposit an arbitrary amount.', 
        'When the required amount is reached, the campaign stops.', 
        'Campaign is considered successful when the author has marked all money transfers.' 
        ]);

INSERT INTO project_types (alias, name, goal_by_amount, goal_by_people, end_by_goal_gain, options) VALUES (
    'money_equal', 'Fair campaign', true, true, false, ARRAY [
        'Partakers agree to split the amount among themselves.',
        'The minimum number of partakers must be recruited.',
        'The number of partakers is not limited.',
        'Fundraising starts on the date specified by the author.',
        'Fair campaign is considered successful when the author has marked all money transfers.' 
        ]);

INSERT INTO project_types (alias, name, goal_by_amount, goal_by_people, end_by_goal_gain, options) VALUES (
    'event_fast', 'Event', false, true, true, ARRAY [
        'Partakers agree to participate in the event.',
        'Event is considered successful when the required number of partakers is reached.'        
        ]);

INSERT INTO project_types (alias, name, goal_by_amount, goal_by_people, end_by_goal_gain, options) VALUES (
    'event_overflow', 'Event+', false, true, false, ARRAY [
        'Partakers agree to participate in the event.',
        'The number of partakers is not limited.',
        'Event+ considered successful if a sufficient number of people have gathered on the event date.' 
        ]);
