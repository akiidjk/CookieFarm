#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <ctype.h>


#define USERNAME_LENGTH 32
#define PASSWORD_LENGTH 16
#define PLAYLIST_NAME_LENGTH 32
#define PLAYLIST_DESCRIPTION_LENGTH 128
#define SONG_NAME_LENGTH 32
#define MAX_SONGS 24
#define MAX_PLAYLISTS 16
#define MAX_USERS 64

#define USER_LINE_LENGTH (USERNAME_LENGTH + PASSWORD_LENGTH + 5)
#define PLAYLIST_LINE_LENGTH (PLAYLIST_NAME_LENGTH + PLAYLIST_DESCRIPTION_LENGTH + (SONG_NAME_LENGTH + 1) * MAX_SONGS + 3)


typedef struct playlist_t
{
	char name[PLAYLIST_NAME_LENGTH];
	char description[PLAYLIST_DESCRIPTION_LENGTH];
	char songs[MAX_SONGS][SONG_NAME_LENGTH];
	uint8_t saved_songs;
} playlist;

typedef struct user_t
{
	char username[USERNAME_LENGTH];
	char password[PASSWORD_LENGTH];
	playlist* playlists[MAX_PLAYLISTS];
	uint8_t saved_playlists;
} user;

user* g_users[MAX_USERS] = { 0 };
user* g_logged = NULL;


user* _username_exists(char*);
uint8_t _check_string(char*);
void _load_users();
void _free_users();

void generate_safe_password(char*, char*);

void user_menu();
void main_menu();

uint8_t register_user();
uint8_t login();

void create_playlist();
void inspect_playlists();
void delete_playlist();


int main()
{
	setvbuf(stdin, NULL, _IONBF, 0);
	setvbuf(stdout, NULL, _IONBF, 0);
	setvbuf(stderr, NULL, _IONBF, 0);

	_load_users();
	main_menu();
	_free_users();

	return 0;
}


user* _username_exists(char* string)
{/*{{{*/
	uint8_t i;
	for (i = 0; i < MAX_USERS && g_users[i] != NULL; ++i)
		if (!strncmp(g_users[i]->username, string, strlen(g_users[i]->username))) break;

	return (i != MAX_USERS && g_users[i] != NULL) ? g_users[i] : NULL;
}/*}}}*/

uint8_t _check_string(char* string)
{/*{{{*/
	string[strcspn(string, "\n")] = '\0';

	while (*string != '\0')
	{
		if (!isalnum(*string) && *string != ' ') break;
		++string;
	}

	return *string == '\0';
}/*}}}*/

void _load_users()
{/*{{{*/
	char user_line[USER_LINE_LENGTH], playlist_line[PLAYLIST_LINE_LENGTH];
	uint8_t num_playlists, i = 0;

	FILE* db = fopen("db.csv", "r"); // TODO: remember permissions (root can write, others can read)
	while (fgets(user_line, USER_LINE_LENGTH, db) != NULL)
	{
		user_line[strcspn(user_line, "\n")] = '\0';

		user* usr = (user*) malloc(sizeof(user));

		strncpy(usr->username, strtok(user_line, ";"), USERNAME_LENGTH);
		strncpy(usr->password, strtok(NULL, ";"), PASSWORD_LENGTH);
		num_playlists = atoi(strtok(NULL, ";"));

		for (uint8_t j = 0; j < num_playlists; ++j)
		{
			if (fgets(playlist_line, PLAYLIST_LINE_LENGTH, db) != NULL)
			{
				playlist_line[strcspn(playlist_line, "\n")] = '\0';

				playlist* pl = (playlist*) malloc(sizeof(playlist));

				strncpy(pl->name, strtok(playlist_line, ";"), PLAYLIST_NAME_LENGTH);
				strncpy(pl->description, strtok(NULL, ";"), PLAYLIST_DESCRIPTION_LENGTH);

				uint8_t k;
				for (k = 0; k < MAX_SONGS; ++k)
				{
					char* tmp = strtok(NULL, ";");
					if (tmp == NULL) break;

					strncpy(pl->songs[k], tmp, SONG_NAME_LENGTH);
				}

				pl->saved_songs = k;
				usr->playlists[j] = pl;
			}
		}

		usr->saved_playlists = num_playlists;
		g_users[i++] = usr;
	}

	fclose(db);
}/*}}}*/

void _free_users()
{/*{{{*/
	char user_line[USER_LINE_LENGTH], playlist_line[PLAYLIST_LINE_LENGTH];

	FILE* db = fopen("db.csv", "w");
	for (uint8_t i = 0; i < MAX_USERS && g_users[i] != NULL; ++i)
	{
		snprintf(user_line, USER_LINE_LENGTH, "%s;%s;%hhu\n", g_users[i]->username, g_users[i]->password, g_users[i]->saved_playlists);
		fwrite(user_line, sizeof(char), strlen(user_line), db);

		for (uint8_t j = 0; j < g_users[i]->saved_playlists; ++j)
		{
			snprintf(playlist_line, PLAYLIST_LINE_LENGTH, "%s;%s", g_users[i]->playlists[j]->name, g_users[i]->playlists[j]->description);
			for (uint8_t k = 0; k < g_users[i]->playlists[j]->saved_songs; ++k)
				snprintf(playlist_line + strlen(playlist_line), SONG_NAME_LENGTH + 1, ";%s", g_users[i]->playlists[j]->songs[k]);
			snprintf(playlist_line + strlen(playlist_line), 2, "\n");

			fwrite(playlist_line, sizeof(char), strlen(playlist_line), db);

			free(g_users[i]->playlists[j]);
			g_users[i]->playlists[j] = NULL;
		}

		free(g_users[i]);
		g_users[i] = NULL;
	}

	fclose(db);
}/*}}}*/

void generate_safe_password(char* username, char* password)
{/*{{{*/
	uint32_t val = 1, tmp;
	for (uint8_t j = 0; j < (uint8_t)(strlen(username) / 4.f) * 4; j += 4)
	{
		tmp = username[j] * (256 * 256 * 256) + username[j + 1] * (256 * 256) + username[j + 2] * 256 + username[j + 3];
		val = (0xDEAD * (val + tmp) + 0xBEEF) & ((1L << 32) - 1);
	}
	snprintf(password, PASSWORD_LENGTH, "%u", val);
}/*}}}*/

void main_menu()
{/*{{{*/
	uint8_t choice, run = 0;

	do
	{
		puts("Vulnify (version: 2.0.1)\n\n");
		puts("Select an option");
		puts("1. Register");
		puts("2. Login");
		puts("0. Exit");
		putchar('>'); putchar(' ');

		choice = getchar(); getchar();
		putchar('\n');

		switch (choice)
		{
			case '1':
				if (!register_user()) run = 1;
				break;
			case '2':
				if (!login()) run = 1;
				break;
			default:
				exit(0);
				break;
		}
	} while (run);

	putchar('\n');

	user_menu();
}/*}}}*/

void user_menu()
{/*{{{*/
	uint8_t choice, run = 0;

	do
	{
		puts("Select an option");
		puts("1. Create a new playlist");
		puts("2. Inspect your playlists");
		puts("3. Delete one of your playlists");
		puts("0. Exit");
		putchar('>'); putchar(' ');

		choice = getchar(); getchar();
		putchar('\n');

		switch (choice)
		{
			case '1':
				create_playlist();
				run = 1;
				break;
			case '2':
				inspect_playlists();
				run = 1;
				break;
			case '3':
				delete_playlist();
				run = 1;
				break;
			default:
				run = 0;
				break;
		}
	} while (run);
}/*}}}*/

uint8_t register_user()
{/*{{{*/
	user* usr = malloc(sizeof(user));
	char* tmp;
	uint8_t i;

	while (1)
	{
		puts("Insert username");
		putchar('>'); putchar(' ');

		tmp = fgets(usr->username, USERNAME_LENGTH, stdin);
		if (tmp == NULL || strlen(tmp) < 1 || !_check_string(tmp) || _username_exists(tmp) != NULL)
		{
			puts("[ERROR] Invalid username");
			continue;
		}

		break;
	}
	usr->username[strcspn(usr->username, "\n")] = '\0';

	puts("Insert password [Empty to generate a safe password automatically]");
	putchar('>'); putchar(' ');

	tmp = fgets(usr->password, PASSWORD_LENGTH, stdin);
	if (tmp == NULL || !_check_string(tmp))
	{
		puts("[ERROR] Invalid password");
		return 0;
	}
	if (tmp[0] == '\0')
	{
		generate_safe_password(usr->username, usr->password);
		printf("Your password is: %s\n", usr->password);
	}
	else usr->password[strcspn(usr->password, "\n")] = '\0';

	for (i = 0; i < MAX_USERS && g_users[i] != NULL; ++i);

	if (i == MAX_USERS)
	{
		puts("[ERROR] Users limit reached");
		return 0;
	}

	puts("User succesfully registered");
	g_logged = g_users[i] = usr;

	return 1;
}/*}}}*/

uint8_t login()
{/*{{{*/
	char username[USERNAME_LENGTH], password[PASSWORD_LENGTH];
	char* tmp;
	user* usr;

	while (1)
	{
		puts("Insert username");
		putchar('>'); putchar(' ');

		tmp = fgets(username, USERNAME_LENGTH, stdin);
		if (tmp == NULL || (usr = _username_exists(username)) == NULL)
		{
			puts("[ERROR] Username not found\n");
			continue;
		}

		break;
	}

	puts("Insert password");
	putchar('>'); putchar(' ');

	tmp = fgets(password, PASSWORD_LENGTH, stdin);
	if (tmp == NULL || strlen(password) < 1 || strncmp(usr->password, password, strlen(password) - 1) != 0)
	{
		puts("[ERROR] Invalid password");
		return 0;
	}

	g_logged = usr;

	return 1;
}/*}}}*/

void create_playlist()
{/*{{{*/
	if (g_logged->saved_playlists == MAX_PLAYLISTS)
	{
		puts("[ERROR] You already reached the maximum number of playlists");
		return;
	}

	playlist* pl = malloc(sizeof(playlist));
	char tmp[SONG_NAME_LENGTH];
	uint8_t num_songs;

	while (1)
	{
		puts("Insert playlist name");
		putchar('>'); putchar(' ');

		if (!fgets(pl->name, PLAYLIST_NAME_LENGTH, stdin) || strlen(pl->name) < 1 || !_check_string(pl->name))
		{
			puts("[ERROR] Invalid name");
			continue;
		}

		break;
	}
	pl->name[strcspn(pl->name, "\n")] = '\0';

	while (1)
	{
		puts("Insert description");
		putchar('>'); putchar(' ');

		if (!fgets(pl->description, PLAYLIST_NAME_LENGTH, stdin) || !_check_string(pl->description))
		{
			puts("[ERROR] Invalid description");
			continue;
		}

		break;
	}
	pl->description[strcspn(pl->description, "\n")] = '\0';

	while (1)
	{
		puts("How many songs do you want to insert?");
		putchar('>'); putchar(' ');

		scanf("%hhu%*c", &num_songs);
	 	if (num_songs >= MAX_SONGS || num_songs == 0)
		{
			printf("[ERROR] Select a positive number smaller than %hhd\n", MAX_SONGS);
			continue;
		}

		break;
	}

	puts("Insert songs");
	uint8_t i = 0;
	while (i < num_songs)
	{
		printf("Song %hhd: ", i);

		if (!fgets(tmp, SONG_NAME_LENGTH, stdin) || !_check_string(tmp))
		{
			puts("[ERROR] Invalid song name");
			continue;
		}
		tmp[strcspn(tmp, "\n")] = '\0';

		strcpy(pl->songs[i++], tmp);
	}
	putchar('\n');

	pl->saved_songs = num_songs;

	g_logged->playlists[g_logged->saved_playlists] = pl;
	g_logged->saved_playlists++;

	puts("Playlist successfully created");
}/*}}}*/

void inspect_playlists()
{/*{{{*/
	for (uint8_t i = 0; i < g_logged->saved_playlists; ++i)
	{
		printf("Playlist %hhd: %s\n\"%s\"\n", i, g_logged->playlists[i]->name, g_logged->playlists[i]->description);
		for (uint8_t j = 0; j < g_logged->playlists[i]->saved_songs; ++j)
			printf("\tSong%hhd: %s\n", j, g_logged->playlists[i]->songs[j]);
		putchar('\n');
	}
	putchar('\n');
}/*}}}*/

void delete_playlist()
{/*{{{*/
	if (g_logged->saved_playlists == 0)
	{
		puts("[ERROR] You have no playlists");
		return;
	}

	puts("Which playlist do you want do delete (use index)?");
	inspect_playlists();

	uint8_t choice;
	while (1)
	{
		putchar('>'); putchar(' ');
		scanf("%hhd%*c", &choice);
		if (choice >= g_logged->saved_playlists)
		{
			puts("[ERROR] Select a valid index");
			continue;
		}

		break;
	}

	playlist* og = g_logged->playlists[choice];

	g_logged->saved_playlists--;
	for (uint8_t i = choice; i < g_logged->saved_playlists; ++i)
		g_logged->playlists[i] = g_logged->playlists[i+1];

	g_logged->playlists[g_logged->saved_playlists] = NULL;
	free(og);

	puts("Playlist successfully deleted");
}/*}}}*/

// vulns:
// can login with just the prefix of a password
// weak password generator
