Projekt zaliczeniowy na przedmiot Programowanie Sieciowe
---

Projekt składa się z aplikacji klienta i aplikacji serwera. Komunikują się one ze sobą za pomocą protokołu TCP wysyłając do siebie odpowiednio spreparowane wiadomości.

Serwer obsługuje żądania `TIME` oraz `DISCONNECT`, które odpowiednio powodują odpowiedź zawierającą aktualny unixowy timestamp z dokładnoscią do milisekund, lub powodują eleganckie zakończenie połączenia klienta z serwerem.

Dodatkowo, jeśli do serwera trafi nieobsługiwane żądanie, serwer odpowie wiadomością `ERROR`.

Zarówno serwer jak i klient zostały napisane w języku Go, bez użycia dodatkowych bibliotek.