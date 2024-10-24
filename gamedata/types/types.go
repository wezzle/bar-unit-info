package types

type (
	UnitPropertiesByRef = map[string]*UnitProperties
	UnitRef             = string
	GridCol             []UnitRef
	GridRow             []GridCol
	Group               []GridRow
	Constructor         = UnitRef
	UnitGrid            map[Constructor]Group
	Lab                 = UnitRef
	LabGrid             map[Lab]GridRow
	WeaponType          = string
	Damage              map[string]float64
	ScarIndices         struct{}
	Shield              struct {
		Repulser         bool
		Smart            bool
		Exterior         bool
		MaxSpeed         float64
		Force            float64
		Radius           float64
		Power            float64
		StartingPower    float64
		PowerRegen       float64
		PowerRegenEnergy float64
		EnergyUse        float64
	}
	Weapon struct {
		BadTargetCategory   []string
		Def                 string
		OnlyTargetCategory  string
		FastAutoRetargeting bool
		MaxAngleDif         int64
		MainDir             string
	}
	WeaponDef struct {
		Name                     string
		WeaponType               WeaponType
		Id                       int64
		CustomParams             map[string]string
		AvoidFriendly            bool
		AvoidFeature             bool
		AvoidNeutral             bool
		AvoidGround              bool
		AvoidCloaked             bool
		CollideEnemy             bool
		CollideFriendly          bool
		CollideFeature           bool
		CollideNeutral           bool
		CollideFireBase          bool
		CollideNonTarget         bool
		CollideGround            bool
		CollideCloaked           bool
		Damage                   Damage
		ExplosionSpeed           float64
		ImpactOnly               bool
		NoSelfDamage             bool
		NoExplode                bool
		Burnblow                 bool
		DamageAreaOfEffect       float64
		EdgeEffectiveness        float64
		CollisionSize            float64
		WeaponVelocity           float64
		StartVelocity            float64
		Weaponacceleration       float64
		ReloadTime               float64
		BurstRate                float64
		Burst                    int64
		Projectiles              int64
		WaterBounce              bool
		GroundBounce             bool
		BounceSlip               float64
		BounceRebound            float64
		NumBounce                int64
		ImpulseFactor            float64
		ImpulseBoost             float64
		CraterMult               float64
		CraterBoost              float64
		CraterAreaOfEffect       float64
		Waterweapon              bool
		Submissile               bool
		FireSubmersed            bool
		Commandfire              bool
		Range                    float64
		Heightmod                float64
		TargetBorder             float64
		CylinderTargeting        float64
		Turret                   bool
		FixedLauncher            bool
		Tolerance                float64
		Firetolerance            float64
		HighTrajectory           int64
		TrajectoryHeight         float64
		Tracks                   bool
		Wobble                   float64
		Dance                    float64
		GravityAffected          bool
		MyGravity                float64
		CanAttackGround          bool
		WeaponTimer              float64
		Flighttime               float64
		Turnrate                 float64
		HeightBoostFactor        float64
		ProximityPriority        float64
		AllowNonBlockingAim      bool
		Accuracy                 float64
		SprayAngle               float64
		MovingAccuracy           float64
		TargetMoveError          float64
		LeadLimit                float64
		LeadBonus                float64
		PredictBoost             float64
		OwnerExpAccWeight        float64
		MinIntensity             float64
		Duration                 float64
		Beamtime                 float64
		Beamburst                bool
		BeamTTL                  int64
		SweepFire                bool
		LargeBeamLaser           bool
		SizeGrowth               float64
		FlameGfxTime             float64
		MetalPerShot             float64
		EnergyPerShot            float64
		FireStarter              float64
		Paralyzer                bool
		ParalyzeTime             int64
		Stockpile                bool
		StockpileTime            float64
		Targetable               int64
		Interceptor              int64
		InterceptedByShieldType  int64
		Coverage                 float64
		InterceptSolo            bool
		DynDamageInverted        bool
		DynDamageExp             float64
		DynDamageMin             float64
		DynDamageRange           float64
		Shield                   Shield
		RechargeDelay            float64
		Model                    string
		Size                     float64
		ScarGlowColorMap         string
		ScarIndices              ScarIndices
		ExplosionScar            bool
		ScarDiameter             float64
		ScarAlpha                float64
		ScarGlow                 float64
		ScarTtl                  float64
		ScarGlowTtl              float64
		ScarDotElimination       float64
		ScarProjVector           [4]float64
		ScarColorTint            [4]float64
		AlwaysVisible            bool
		CameraShake              float64
		SmokeTrail               bool
		SmokeTrailCastShadow     bool
		SmokePeriod              int64
		SmokeTime                int64
		SmokeSize                float64
		SmokeColor               float64
		CastShadow               bool
		SizeDecay                float64
		AlphaDecay               float64
		Separation               float64
		NoGap                    bool
		Stages                   int64
		LodDistance              int64
		Thickness                float64
		CoreThickness            float64
		LaserFlareSize           float64
		TileLength               float64
		ScrollSpeed              float64
		PulseSpeed               float64
		BeamDecay                float64
		FalloffRate              float64
		Hardstop                 bool
		RgbColor                 [3]float64
		RgbColor2                [3]float64
		Intensity                float64
		Colormap                 string
		CegTag                   string
		ExplosionGenerator       string
		BounceExplosionGenerator string
		SoundTrigger             bool
		SoundStart               string
		SoundHitDry              string
		SoundHitWet              string
		SoundStartVolume         float64
		SoundHitDryVolume        float64
		SoundHitWetVolume        float64
	}
	CustomParams struct {
		TechLevel int
		UnitGroup string
	}
	UnitProperties struct {
		Ref            UnitRef
		MetalCost      int64
		EnergyCost     int64
		Buildtime      int64
		BuildOptions   []UnitRef
		Health         int64
		SightDistance  int64
		Speed          float64
		Buildpower     int64
		SonarDistance  int64
		RadarDistance  int64
		JammerDistance int64
		WeaponDefs     map[string]WeaponDef
		Weapons        []Weapon
		CustomParams   CustomParams
	}
	Translations struct {
		Units struct {
			Factions                  map[string]string  `json:"factions"`
			Dead                      string             `json:"dead"`
			Heap                      string             `json:"heap"`
			DecoyCommanderNameTag     string             `json:"decoyCommanderNameTag"`
			Scavenger                 string             `json:"scavenger"`
			ScavCommanderNameTag      string             `json:"scavCommanderNameTag"`
			ScavDecoyCommanderNameTag string             `json:"scavDecoyCommanderNameTag"`
			Names                     map[UnitRef]string `json:"names"`
			Descriptions              map[UnitRef]string `json:"descriptions"`
		} `json:"units"`
	}
)
