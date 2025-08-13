# Berechnung Werte und Boni

        pa: 1d100 + 4xIn/10 - 20,               // Persönliche Ausstrahlung
        wk: 1d100 + 2x(Ko/10 + In/10) - 20,       // Willenskraft
        lp_max: 1d3 + 7 + (Ko/10),               // Lebenspunkte (Maximum)
                Gnome -3, Halblinge -2, Zwerge +1
        ausdauer_bonus: Ko/10 + St/20, // Ausdauer Bonus
        ap_max: 1d3 + 1 + ausdauer_bonus,
                Barbar, Krieger, Waldläufer +2
                andere Kämpfer, Schamane +1
        schadens_bonus: St/20 + Gs/30 -3, // Schadens Bonus
        b_max: 4d3 +16,                         // Bewegungsweite
                Gnome: 2d3 +8
                Halblinge: 2d3 +8
                Zwerge: 3d3 + 12
        angriffs_bonus:                      // Angriffs Bonus
                GS 01-05: -2
                GS 06-20: -1
                GS 21-80: 0
                GS 81-95: +1
                GS 96-100: +2
        abwehr_bonus:                       // Abwehr Bonus
                GW 01-05: -2
                GW 06-20: -1
                GW 21-80: 0
                GW 81-95: +1
                GW 96-100: +2
        zauber_bonus:                      // Zauber Bonus
                Zt 01-05: -2
                Zt 06-20: -1
                Zt 21-80: 0
                Zt 81-95: +1
                Zt 96-100: +2
        resistenz_bonus_koerper:         // Resistenz Bonus Körper
                für Menschen:
                        Ko 01-05: -2
                        Ko 06-20: -1
                        Ko 21-80: 0
                        Ko 81-95: +1
                        Ko 96-100: +2
                        zusätzlich: Kämpfer +1, Zauberer +2
                für Nicht Menschen:
                        Elfen: +2
                        Gnome: +4
                        Halblinge: +4
                        Zwerge: +3
                        zusätzlich Kämpfer +1, Zauberer +2
        resistenz_bonus_geist:          // Resistenz Bonus Geist
                für Menschen:
                        In 01-05: -2
                        In 06-20: -1
                        In 21-80: 0
                        In 81-95: +1
                        In 96-100: +2
                        zusätzlich: Zauberer +2
                für Nicht Menschen:
                        Elfen: +2
                        Gnome: +4
                        Halblinge: +4
                        Zwerge: +3
                        zusätzlich Kämpfer 0, Zauberer +2
        resistenz_koerper: 11 + resistenz_bonus_koerper, // Resistenz Körper
        resistenz_geist: 11 + resistenz_bonus_geist, // Resistenz Geist
        abwehr: 11 + abwehr_bonus, // Abwehr
        zaubern: 11 + zauber_bonus, // Zaubern
        raufen: (St + GW)/20 + angriffs_bonus             // Raufen
                zusätzlich: Zwerge +1
        Fertigkeit:
                Trinken: Ko/10
                Wahrnehmung: +6
                Sprache: +12    (Muttersprache)
                Sprache: +12    (Gemeinsprache)